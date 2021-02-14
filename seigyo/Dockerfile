FROM python:3.7-buster

### mine ###

RUN curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py \
    && python3 get-pip.py \
    && pip3 install grpcio-tools

RUN pip install urllib3 pyyaml

COPY . /root

EXPOSE 80 8081 9000

WORKDIR /root
CMD ["/root/start.sh"]
