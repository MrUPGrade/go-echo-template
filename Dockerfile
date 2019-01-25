FROM debian

ADD echoapi /echoapi
RUN chmod +x /echoapi

CMD ["/echoapi"]