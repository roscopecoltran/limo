#!/usr/bin/env python3

# Copyright (c) 2016 Salle, Alexandre <atsalle@inf.ufrgs.br>
# Author: Salle, Alexandre <atsalle@inf.ufrgs.br>
# 
# Permission is hereby granted, free of charge, to any person obtaining a copy of
# this software and associated documentation files (the "Software"), to deal in
# the Software without restriction, including without limitation the rights to
# use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
# the Software, and to permit persons to whom the Software is furnished to do so,
# subject to the following conditions:
# 
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
# 
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
# FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
# COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
# IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
# CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

from __future__ import print_function
import argparse
import subprocess
import os
import sys

arg_map = [
	('train', 'corpus', '', True, 'tokenized corpus used for training'),
	('output', 'output', '', True, 'file where to output vectors'),
	('size', 'dim', '100', False, 'dimensions of vectors'),
	('window', 'window', '2', False, 'window to each side of target word'),
	('sample', 'subsample', '1e-3', False, 'subsampling corpus'),
	('negative', 'negative', '5', False, '# of negative samples'),
	('threads', 'threads', '12', False, '# of threads'),
	('iter', 'iterations', '5', False, '# of iterations'),
	('min-count', 'minfreq', '5', False, 'remove from vocabulary words that appear less than # times'),
	('alpha', 'alpha', '0.025', False, 'learning rate'),
	('debug', 'verbose', '2', False, 'debug info (0 for none, 1 for some, 2 for lots)'),
	('save-vocab', 'savevocab', None, False, 'save vocab to file'),
	('read-vocab', 'readvocab', None, False, 'read vocab from file'),
]

LEXVEC = os.environ.get('LEXVEC', './lexvec')

ignore = ['hs', 'classes', 'binary', 'cbow']

parser = argparse.ArgumentParser(description="word2vec interface to lexvec", formatter_class=argparse.ArgumentDefaultsHelpFormatter)
for arg in arg_map:
	parser.add_argument('-' + arg[0], dest=arg[1], default=arg[2], required=arg[3], help=arg[4])
for arg in ignore:
	parser.add_argument('-' + arg, dest=arg, default=None)
args = parser.parse_args()

p = vars(args)
output = p['output']
del p['output']

args = map(lambda x: ('-' + x[0], x[1]), filter(lambda x: x[0] not in ignore and x[1], p.items()))
args = [y for x in args for y in x]

cmd = [LEXVEC] + args
print(' '.join(cmd))
subprocess.Popen(cmd, stdout=open(output, 'w'), stderr=sys.stderr).wait()
