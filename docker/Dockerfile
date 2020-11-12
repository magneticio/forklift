FROM alpine:3.12

RUN apk add --no-cache \
    ca-certificates \
    curl \
    nano

RUN \
  version=$(curl -s https://api.github.com/repos/magneticio/forklift/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/') && \
  base=https://github.com/magneticio/forklift/releases/download/$version && \
  curl -sL $base/forklift-$(uname -s)-$(uname -m) >/usr/local/bin/forklift && \
  chmod +x /usr/local/bin/forklift
  
CMD ["/bin/sh"]
