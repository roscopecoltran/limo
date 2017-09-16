###########################################################################
#		  
#  Build the image:                                               		  
#    $ docker build -t snk-apk-file:go1.9-alpine3.6 --no-cache . 			# longer but more accurate
#    $ docker build -t snk-apk-file:go1.9-alpine3.6 . 						# faster but increase mistakes
#                                                                 		  
#  Run the container:                                             		  
#    $ docker run -it --rm -v $(pwd)/shared:/shared snk-apk-file:go1.9-alpine3.6
#    $ docker run -d --name snk-apk-file -v $(pwd)/shared:/shared snk-apk-file:go1.9-alpine3.6
#                                                              		  
###########################################################################

## LEVEL1 ###############################################################################################################


FROM alpine:3.6
LABEL maintainer "Luc Michalski <michalski.luc@gmail.com>"

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN	apk add --no-cache \
	ca-certificates

COPY . /go/src/github.com/jessfraz/apk-file

RUN set -x \
	&& apk add --no-cache --virtual .build-deps \
		go \
		git \
		gcc \
		libc-dev \
		libgcc \
	&& cd /go/src/github.com/jessfraz/apk-file \
	&& go build -o /usr/bin/apk-file . \
	&& apk del .build-deps \
	&& rm -rf /go \
	&& echo "Build complete."


ENTRYPOINT [ "apk-file" ]