godoc-proxy is a go web server listening on 127.0.0.1:80. All communication to this web
server is sent to godoc.org with the exception of /-/site.css.

Customize override.css to your preferences. override.css currently changes
colors and allows the full page width.

### Install
- go generate
- go install
- use godoc-proxy.service for Linux or godoc.plist for launchd

#### License 

Use of this source code is governed by a BSD-2-Clause license that can be found
in the LICENSE file.

[![BSD-2-Clause License](img/osi_logo_100X133_90ppi_0.png)](https://opensource.org/)
