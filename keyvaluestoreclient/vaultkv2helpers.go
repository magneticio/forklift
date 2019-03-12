package keyvaluestoreclient

import (
	"encoding/json"
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	kvbuilder "github.com/hashicorp/vault/helper/kv-builder"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/kr/text"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/ryanuber/columnize"
)

const (
	// maxLineLength is the maximum width of any line.
	maxLineLength int = 78

	// notSetValue is a flag value for a not-set value
	notSetValue = "(not set)"
)

// extractListData reads the secret and returns a typed list of data and a
// boolean indicating whether the extraction was successful.
func extractListData(secret *api.Secret) ([]interface{}, bool) {
	if secret == nil || secret.Data == nil {
		return nil, false
	}

	k, ok := secret.Data["keys"]
	if !ok || k == nil {
		return nil, false
	}

	i, ok := k.([]interface{})
	return i, ok
}

// sanitizePath removes any leading or trailing things from a "path".
func sanitizePath(s string) string {
	return ensureNoTrailingSlash(ensureNoLeadingSlash(strings.TrimSpace(s)))
}

// ensureTrailingSlash ensures the given string has a trailing slash.
func ensureTrailingSlash(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	for len(s) > 0 && s[len(s)-1] != '/' {
		s = s + "/"
	}
	return s
}

// ensureNoTrailingSlash ensures the given string has a trailing slash.
func ensureNoTrailingSlash(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	for len(s) > 0 && s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}
	return s
}

// ensureNoLeadingSlash ensures the given string has a trailing slash.
func ensureNoLeadingSlash(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	for len(s) > 0 && s[0] == '/' {
		s = s[1:]
	}
	return s
}

// columnOuput prints the list of items as a table with no headers.
func columnOutput(list []string, c *columnize.Config) string {
	if len(list) == 0 {
		return ""
	}

	if c == nil {
		c = &columnize.Config{}
	}
	if c.Glue == "" {
		c.Glue = "    "
	}
	if c.Empty == "" {
		c.Empty = "n/a"
	}

	return columnize.Format(list, c)
}

// tableOutput prints the list of items as columns, where the first row is
// the list of headers.
func tableOutput(list []string, c *columnize.Config) string {
	if len(list) == 0 {
		return ""
	}

	delim := "|"
	if c != nil && c.Delim != "" {
		delim = c.Delim
	}

	underline := ""
	headers := strings.Split(list[0], delim)
	for i, h := range headers {
		h = strings.TrimSpace(h)
		u := strings.Repeat("-", len(h))

		underline = underline + u
		if i != len(headers)-1 {
			underline = underline + delim
		}
	}

	list = append(list, "")
	copy(list[2:], list[1:])
	list[1] = underline

	return columnOutput(list, c)
}

// parseArgsData parses the given args in the format key=value into a map of
// the provided arguments. The given reader can also supply key=value pairs.
func parseArgsData(stdin io.Reader, args []string) (map[string]interface{}, error) {
	builder := &kvbuilder.Builder{Stdin: stdin}
	if err := builder.Add(args...); err != nil {
		return nil, err
	}

	return builder.Map(), nil
}

// parseArgsDataString parses the args data and returns the values as strings.
// If the values cannot be represented as strings, an error is returned.
func parseArgsDataString(stdin io.Reader, args []string) (map[string]string, error) {
	raw, err := parseArgsData(stdin, args)
	if err != nil {
		return nil, err
	}

	var result map[string]string
	if err := mapstructure.WeakDecode(raw, &result); err != nil {
		return nil, errors.Wrap(err, "failed to convert values to strings")
	}
	if result == nil {
		result = make(map[string]string)
	}
	return result, nil
}

// parseArgsDataStringLists parses the args data and returns the values as
// string lists. If the values cannot be represented as strings, an error is
// returned.
func parseArgsDataStringLists(stdin io.Reader, args []string) (map[string][]string, error) {
	raw, err := parseArgsData(stdin, args)
	if err != nil {
		return nil, err
	}

	var result map[string][]string
	if err := mapstructure.WeakDecode(raw, &result); err != nil {
		return nil, errors.Wrap(err, "failed to convert values to strings")
	}
	return result, nil
}

// truncateToSeconds truncates the given duration to the number of seconds. If
// the duration is less than 1s, it is returned as 0. The integer represents
// the whole number unit of seconds for the duration.
func truncateToSeconds(d time.Duration) int {
	d = d.Truncate(1 * time.Second)

	// Handle the case where someone requested a ridiculously short increment -
	// increments must be larger than a second.
	if d < 1*time.Second {
		return 0
	}

	return int(d.Seconds())
}

// printKeyStatus prints the KeyStatus response from the API.
func printKeyStatus(ks *api.KeyStatus) string {
	return columnOutput([]string{
		fmt.Sprintf("Key Term | %d", ks.Term),
		fmt.Sprintf("Install Time | %s", ks.InstallTime.UTC().Format(time.RFC822)),
	}, nil)
}

// expandPath takes a filepath and returns the full expanded path, accounting
// for user-relative things like ~/.
func expandPath(s string) string {
	if s == "" {
		return ""
	}

	e, err := homedir.Expand(s)
	if err != nil {
		return s
	}
	return e
}

// wrapAtLengthWithPadding wraps the given text at the maxLineLength, taking
// into account any provided left padding.
func wrapAtLengthWithPadding(s string, pad int) string {
	wrapped := text.Wrap(s, maxLineLength-pad)
	lines := strings.Split(wrapped, "\n")
	for i, line := range lines {
		lines[i] = strings.Repeat(" ", pad) + line
	}
	return strings.Join(lines, "\n")
}

// wrapAtLength wraps the given text to maxLineLength.
func wrapAtLength(s string) string {
	return wrapAtLengthWithPadding(s, 0)
}

// ttlToAPI converts a user-supplied ttl into an API-compatible string. If
// the TTL is 0, this returns the empty string. If the TTL is negative, this
// returns "system" to indicate to use the system values. Otherwise, the
// time.Duration ttl is used.
func ttlToAPI(d time.Duration) string {
	if d == 0 {
		return ""
	}

	if d < 0 {
		return "system"
	}

	return d.String()
}

// humanDuration prints the time duration without those pesky zeros.
func humanDuration(d time.Duration) string {
	if d == 0 {
		return "0s"
	}

	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if idx := strings.Index(s, "h0m"); idx > 0 {
		s = s[:idx+1] + s[idx+3:]
	}
	return s
}

// humanDurationInt prints the given int as if it were a time.Duration  number
// of seconds.
func humanDurationInt(i interface{}) interface{} {
	switch i.(type) {
	case int:
		return humanDuration(time.Duration(i.(int)) * time.Second)
	case int64:
		return humanDuration(time.Duration(i.(int64)) * time.Second)
	case json.Number:
		if i, err := i.(json.Number).Int64(); err == nil {
			return humanDuration(time.Duration(i) * time.Second)
		}
	}

	// If we don't know what type it is, just return the original value
	return i
}

func kvReadRequest(client *api.Client, path string, params map[string]string) (*api.Secret, error) {
	r := client.NewRequest("GET", "/v1/"+path)
	for k, v := range params {
		r.Params.Set(k, v)
	}
	resp, err := client.RawRequest(r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp != nil && resp.StatusCode == 404 {
		secret, parseErr := api.ParseSecret(resp.Body)
		switch parseErr {
		case nil:
		case io.EOF:
			return nil, nil
		default:
			return nil, err
		}
		if secret != nil && (len(secret.Warnings) > 0 || len(secret.Data) > 0) {
			return secret, nil
		}
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return api.ParseSecret(resp.Body)
}

func kvPreflightVersionRequest(client *api.Client, path string) (string, int, error) {
	// We don't want to use a wrapping call here so save any custom value and
	// restore after
	currentWrappingLookupFunc := client.CurrentWrappingLookupFunc()
	client.SetWrappingLookupFunc(nil)
	defer client.SetWrappingLookupFunc(currentWrappingLookupFunc)
	currentOutputCurlString := client.OutputCurlString()
	client.SetOutputCurlString(false)
	defer client.SetOutputCurlString(currentOutputCurlString)

	r := client.NewRequest("GET", "/v1/sys/internal/ui/mounts/"+path)
	resp, err := client.RawRequest(r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		// If we get a 404 we are using an older version of vault, default to
		// version 1
		if resp != nil && resp.StatusCode == 404 {
			return "", 1, nil
		}

		return "", 0, err
	}

	secret, err := api.ParseSecret(resp.Body)
	if err != nil {
		return "", 0, err
	}
	var mountPath string
	if mountPathRaw, ok := secret.Data["path"]; ok {
		mountPath = mountPathRaw.(string)
	}
	options := secret.Data["options"]
	if options == nil {
		return mountPath, 1, nil
	}
	versionRaw := options.(map[string]interface{})["version"]
	if versionRaw == nil {
		return mountPath, 1, nil
	}
	version := versionRaw.(string)
	switch version {
	case "", "1":
		return mountPath, 1, nil
	case "2":
		return mountPath, 2, nil
	}

	return mountPath, 1, nil
}

func isKVv2(path string, client *api.Client) (string, bool, error) {
	mountPath, version, err := kvPreflightVersionRequest(client, path)
	if err != nil {
		return "", false, err
	}

	return mountPath, version == 2, nil
}

func addPrefixToVKVPath(p, mountPath, apiPrefix string) string {
	switch {
	case p == mountPath, p == strings.TrimSuffix(mountPath, "/"):
		return path.Join(mountPath, apiPrefix)
	default:
		p = strings.TrimPrefix(p, mountPath)
		return path.Join(mountPath, apiPrefix, p)
	}
}

func getHeaderForMap(header string, data map[string]interface{}) string {
	maxKey := 0
	for k := range data {
		if len(k) > maxKey {
			maxKey = len(k)
		}
	}

	// 4 for the column spaces and 5 for the len("value")
	totalLen := maxKey + 4 + 5

	equalSigns := totalLen - (len(header) + 2)

	// If we have zero or fewer equal signs bump it back up to two on either
	// side of the header.
	if equalSigns <= 0 {
		equalSigns = 4
	}

	// If the number of equal signs is not divisible by two add a sign.
	if equalSigns%2 != 0 {
		equalSigns = equalSigns + 1
	}

	return fmt.Sprintf("%s %s %s", strings.Repeat("=", equalSigns/2), header, strings.Repeat("=", equalSigns/2))
}

func kvParseVersionsFlags(versions []string) []string {
	versionsOut := make([]string, 0, len(versions))
	for _, v := range versions {
		versionsOut = append(versionsOut, strutil.ParseStringSlice(v, ",")...)
	}

	return versionsOut
}
