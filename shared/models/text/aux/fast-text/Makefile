# determine platform
ifeq (Darwin, $(findstring Darwin, $(shell uname -a)))
  PLATFORM 			:= macosx
  GO_BUILD_OS 		:= darwin
else
  PLATFORM 			:= Linux
  GO_BUILD_OS 		:= linux
endif

PGET_SPLIT_PROCESSOR_COUNT 			:= 6

ARCHIVE_ENGINE_FILE_EXT 			:= zip
ARCHIVE_CHUNKS_FILE_SIZE_MB			:= 60

FASTTEXT_VECTORS_REMOTE_URL 		:= https://s3-us-west-1.amazonaws.com/fasttext-vectors/wiki.en.zip
FASTTEXT_VECTORS_REMOTE_CHECKSUM	:= sha1:e74820b00ecd26c9d4e716d5cec9ac53a75b0eea

FASTTEXT_VECTORS_BASENAME 			:= fasttext-vectors-wiki.en
FASTTEXT_VECTORS_FILEPATH 			:= $(CURDIR)/$(FASTTEXT_VECTORS_BASENAME)

FASTTEXT_VECTORS_ARCHIVE_FILENAME 	:= $(FASTTEXT_VECTORS_BASENAME).$(ARCHIVE_ENGINE_FILE_EXT)
FASTTEXT_VECTORS_ARCHIVE_FILEPATH 	:= $(CURDIR)/$(FASTTEXT_VECTORS_ARCHIVE_FILENAME).$(ARCHIVE_ENGINE_FILE_EXT)

download: download-pget

download-pget:
	pget -p $(PGET_SPLIT_PROCESSOR_COUNT) --target-dir $(CURDIR) --output ${FASTTEXT_VECTORS_ARCHIVE_FILENAME} $(FASTTEXT_VECTORS_REMOTE_URL)
	# cat $(FASTTEXT_VECTORS_ARCHIVE_FILEPATH) | checksum sha1
	unzip ${FASTTEXT_VECTORS_ARCHIVE_FILEPATH}

extract:
	@if [ -f $(GOOGLE_NEWS_ARCHIVE_FILEPATH) ]; then gunzip $(GOOGLE_NEWS_ARCHIVE_FILEPATH); fi

compress: extract
	@if [ ! -f $(GOOGLE_NEWS_ARCHIVE_FILEPATH) ]; then tar zcvf $(GOOGLE_NEWS_ARCHIVE_FILEPATH) $(GOOGLE_NEWS_ARCHIVE_BASENAME); fi

compress-split: compress-split-$(PLATFORM)

compress-split-linux:
	@tar zcvf - $(GOOGLE_NEWS_ARCHIVE_BASENAME) | split --bytes=$(ARCHIVE_CHUNKS_FILE_SIZE_MB)MB - $(GOOGLE_NEWS_ARCHIVE_FILENAME).

compress-split-macosx:
	@tar zcvf - $(GOOGLE_NEWS_ARCHIVE_BASENAME) | split -b $(ARCHIVE_CHUNKS_FILE_SIZE_MB)m - $(GOOGLE_NEWS_ARCHIVE_FILENAME).
