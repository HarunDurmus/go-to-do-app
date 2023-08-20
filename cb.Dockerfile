FROM couchbase:latest

ENV CB_ADMIN admin
ENV CB_ADMIN_PASSWORD 123456
ENV CB_BUCKET todoapp
ENV CB_BUCKET_RAM_SIZE 512
ENV CB_SERVICES "data,index,query,fts"

COPY ./configure-couchbase.sh /


RUN ["chmod", "+x", "/configure-couchbase.sh"]

ENTRYPOINT ["/configure-couchbase.sh"]
