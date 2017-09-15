
# determine platform
ifeq (Darwin, $(findstring Darwin, $(shell uname -a)))
  PLATFORM 			:= macosx
  GO_BUILD_OS 		:= darwin
else
  PLATFORM 			:= Linux
  GO_BUILD_OS 		:= linux
endif

ARCHIVE_ENGINE_FILE_EXT 		:= gz
ARCHIVE_CHUNKS_FILE_SIZE_MB		:= 60

GOOGLE_NEWS_ARCHIVE_BASENAME 	:= google-news-vectors-negative300.bin
GOOGLE_NEWS_ARCHIVE_FILENAME 	:= $(GOOGLE_NEWS_ARCHIVE_BASENAME).$(ARCHIVE_ENGINE_FILE_EXT)
GOOGLE_NEWS_ARCHIVE_FILEPATH 	:= $(CURDIR)/$(GOOGLE_NEWS_ARCHIVE_BASENAME).$(ARCHIVE_ENGINE_FILE_EXT)

extract:
	@if [ -f $(GOOGLE_NEWS_ARCHIVE_FILEPATH) ]; then gunzip $(GOOGLE_NEWS_ARCHIVE_FILEPATH); fi

compress: extract
	@if [ ! -f $(GOOGLE_NEWS_ARCHIVE_FILEPATH) ]; then tar zcvf $(GOOGLE_NEWS_ARCHIVE_FILEPATH) $(GOOGLE_NEWS_ARCHIVE_BASENAME); fi

compress-split: compress-split-$(PLATFORM)

compress-split-linux:
	@tar zcvf - $(GOOGLE_NEWS_ARCHIVE_BASENAME) | split --bytes=$(ARCHIVE_CHUNKS_FILE_SIZE_MB)MB - $(GOOGLE_NEWS_ARCHIVE_FILENAME).

compress-split-macosx:
	@tar zcvf - $(GOOGLE_NEWS_ARCHIVE_BASENAME) | split -b $(ARCHIVE_CHUNKS_FILE_SIZE_MB)m - $(GOOGLE_NEWS_ARCHIVE_FILENAME).
