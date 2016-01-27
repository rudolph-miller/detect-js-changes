# Detect JS Changes

detect changes after some operations, for example deploying.

# Install

```sh
$ go get github.com/Rudolph-Miller/detect-js-changes
```

# Config

```
$ cat detect_config.yml
default:
  urls:
  - https://raw.githubusercontent.com/lodash/lodash/4.0.1/dist/lodash.core.min.js
  - https://raw.githubusercontent.com/lodash/lodash/4.0.0/dist/lodash.core.min.js
  tmp_dir: ./tmp
  ignore_keywords:
  - sample keyword
development:
  ignore_keywords:
  - hello world
```

# Usage

```
$ detect-js-changes reset
Reset: tmp/detect_js_changes_download_1
Reset: tmp/detect_js_changes_download_2

$ detect-js-changes download
Directory: tmp/detect_js_changes_download_1
Download: https://raw.githubusercontent.com/lodash/lodash/4.0.1/dist/lodash.core.min.js as file_0
Download: https://raw.githubusercontent.com/lodash/lodash/4.0.0/dist/lodash.core.min.js as file_1

$ deploy
...

$ detect-js-changes download
Directory: tmp/detect_js_changes_download_2
Download: https://raw.githubusercontent.com/lodash/lodash/4.0.1/dist/lodash.core.min.js as file_0
Download: https://raw.githubusercontent.com/lodash/lodash/4.0.0/dist/lodash.core.min.js as file_1

$ detect-js-changes detect
Detecting: https://raw.githubusercontent.com/lodash/lodash/4.0.1/dist/lodash.core.min.js
Result: https://raw.githubusercontent.com/lodash/lodash/4.0.1/dist/lodash.core.min.js has no changes
Detecting: https://raw.githubusercontent.com/lodash/lodash/4.0.0/dist/lodash.core.min.js
Result: https://raw.githubusercontent.com/lodash/lodash/4.0.0/dist/lodash.core.min.js has no changes
```

# env

```
$ detect-js-changes --env development detect
$ detect-js-changes -e development detect
$ ENV=development detect-js-changes detect
```

# reset

```
$ detect-js-changes reset
Reset: tmp/detect_js_changes_download_1
Reset: tmp/detect_js_changes_download_2

$ detect-js-changes reset 1
Reset: tmp/detect_js_changes_download_1

$ detect-js-changes reset 2
Reset: tmp/detect_js_changes_download_2
```

# download

```
$ detect-js-changes download
```

# detect

```
$ detect-js-changes detect
```

Result has 3 ways below.

```
Result: URL has some changes
Result: URL has no changes
Result: URL has ignored changes
```
