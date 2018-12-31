FROM scratch
EXPOSE 8080
ENTRYPOINT ["/cassuservice"]
COPY ./bin/ /