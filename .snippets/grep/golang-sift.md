sift --multiline --filename --smart-case 'build\(\)\s{(.*?)\n}' --replace '$1' --write-conf

