FROM ubuntu:22.04

RUN mkdir /service && apt-get update && apt-get install socat -y && apt-get install g++ -y && apt-get install make -y


WORKDIR /service

COPY ./ ./

RUN useradd sfuser

RUN chown -R sfuser /service 

USER sfuser

RUN mkdir /service/storage

EXPOSE 7490

CMD ["./run.sh"]
