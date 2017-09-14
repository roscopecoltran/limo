#!/usr/bin/env python3

import codecs
from collections import defaultdict
import gzip
import io
import json
import os
import pickle
import re
import shutil
import subprocess
import sys

from nltk.stem.snowball import SnowballStemmer
import numpy
from pygments import highlight
from pygments.lexers import get_lexer_for_filename
from pygments.formatter import Formatter

NAME_BREAKUP_RE = re.compile(r"[^a-zA-Z]+")

def extract_names(t):
    t = t.strip()
    prev_p = [""]

    def ret(name):
        r = name.lower()
        if len(name) >= 3:
            yield r
            if prev_p[0]:
                yield prev_p[0] + r
                prev_p[0] = ""
        else:
            prev_p[0] = r

    for part in NAME_BREAKUP_RE.split(t):
        if not part:
            continue
        prev = t[0]
        pos = 0
        for i in range(1, len(part)):
            this = part[i]
            if prev.islower() and this.isupper():
                yield from ret(part[pos:i])
                pos = i
            elif prev.isupper() and this.islower():
                if 0 < i - 1 - pos <= 3:
                    yield from ret(part[pos:i - 1])
                    pos = i - 1
                elif i - 1 > pos:
                    yield from ret(part[pos:i])
                    pos = i
            prev = this
        last = part[pos:]
        if last:
            yield from ret(last)


def out(text):
    stdout = sys.stdout
    stdout.write(" " * 80 + "\r")
    stdout.write(text)
    stdout.write("\r")
    stdout.flush()


def main():
    out("loading the model...")
    with gzip.open("repo_topic_modelling.pickle.gz", "rb") as gzfin:
        with io.BufferedReader(gzfin) as fin:
            model, idf = pickle.load(fin)
    word_index = {w: i for i, w in enumerate(model.index)}
    repo = sys.argv[1]
    out("cloning %s..." % repo)
    if os.path.exists("repo"):
        shutil.rmtree("repo")
    subprocess.check_call(
        ["git", "clone", "--depth=1", "--recursive", "https://github.com/" + repo, "repo"],
        stderr=subprocess.DEVNULL)
    out("classifying the files...")
    linfiles = subprocess.check_output(["linguist", "repo", "--json"])
    files = json.loads(linfiles.decode("utf-8"))
    names = defaultdict(int)
    stemmer = SnowballStemmer(language="english")

    def format_callback(tokensource, _):
        for ttype, value in tokensource:
            if ttype[0] == "Name":
                for name in extract_names(value):
                    names[stemmer.stem(name)] += 1

    formatter = Formatter()
    formatter.format = format_callback

    for lang, lfiles in files.items():
        out("parsing %s files..." % lang)
        lexer = None
        for file in lfiles:
            with codecs.open(os.path.join("repo", file), "r", "utf-8") as fin:
                code = fin.read()
            if lexer is None:
                try:
                    lexer = get_lexer_for_filename(file, code)
                except:
                    continue
            highlight(code, lexer, formatter)

    out("applying the model...")
    vec = numpy.zeros(len(word_index))
    for n, f in names.items():
        i = word_index.get(n)
        if i is not None:
            vec[i] = f
    vec = numpy.log(vec + 1) * idf
    topics = vec.dot(model.values)
    out("")
    for i in topics.argsort()[:-6:-1]:
        print("%.2f\t%s" % (topics[i], model.columns[i]))


if __name__ == "__main__":
    sys.exit(main())
