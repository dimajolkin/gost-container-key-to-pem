FROM sokolko/export-cryptopro-cert as openssl-gost

FROM ubuntu:last as base
COPY --from=openssl-gost /usr/local/ssl /usr/local/ssl
COPY --from=openssl-gost /usr/bin/openssl /usr/bin/openssl
COPY --from=openssl-gost /lib/ld-musl-*.so.1 /lib
COPY --from=openssl-gost /bin/get-cpcert /bin/get-cpcert

RUN mkdir /data

CMD sleep 100