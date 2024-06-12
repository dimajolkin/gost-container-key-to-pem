FROM sokolko/export-cryptopro-cert as openssl-gost

FROM ubuntu:20.04 as base

COPY --from=openssl-gost /usr/local/ssl /usr/local/ssl
COPY --from=openssl-gost /usr/bin/openssl /usr/bin/openssl
COPY --from=openssl-gost /lib/ld-musl-*.so.1 /lib

# openssl ciphers|tr ':' '\n'|grep GOST
ENV OPENSSL_CONF /usr/local/ssl/openssl.cnf

COPY --from=openssl-gost /bin/get-cpcert /bin/get-cpcert

RUN mkdir /data

CMD sleep 100