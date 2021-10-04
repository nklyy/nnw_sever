FROM ubuntu:18.04
RUN apt-get update
RUN apt-get install build-essential autoconf autotools-dev libtool pkg-config libboost-all-dev libssl-dev libprotobuf-dev protobuf-compiler libevent-dev libqt4-dev libcanberra-gtk-module libdb-dev libdb++-dev bsdmainutils -y
WORKDIR /src/bitcoin
COPY ./bitcoin .
RUN ./autogen.sh
RUN ./configure --with-incompatible-bdb
RUN make
RUN make install
RUN mkdir -p /root/bitcoind-simnet/