FROM golang:alpine AS build
#pull dependencies required for downloading and building libndctl
RUN apk add --update build-base autoconf automake libtool libuuid json-c kmod asciidoc xmlto kmod-dev eudev-dev util-linux-dev json-c-dev linux-headers wget tar file

WORKDIR /
RUN wget https://github.com/pmem/ndctl/archive/v62.tar.gz
RUN tar zxvf v62.tar.gz
WORKDIR /ndctl-62
RUN ./autogen.sh
RUN ./configure CFLAGS='-g -O2' --prefix=/usr --sysconfdir=/etc --libdir=/usr/lib --disable-docs with_systemd_unit_dir=no
RUN make install

# build pmem-vgm
ADD . /go/src/github.com/intel/pmem-csi
WORKDIR /go/src/github.com/intel/pmem-csi
RUN make pmem-vgm
RUN mv ./_output/pmem-vgm /go/bin/

# build clean container
FROM golang:alpine
LABEL maintainers="Intel"
LABEL description="Pmem CSI driver volume group manager"
RUN apk add --update kmod eudev util-linux libuuid lvm2 file
# move required binaries and libraries to clean container
COPY --from=build /usr/lib/libndctl* /usr/lib/
COPY --from=build /usr/lib/libdaxctl* /usr/lib/
RUN mkdir -p /go/bin
COPY --from=build /go/bin/pmem-vgm /go/bin/

ENTRYPOINT ["/go/bin/pmem-vgm"]
