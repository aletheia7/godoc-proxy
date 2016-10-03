godoc-proxy is a go web server listening on :80. All communication to this web
server is sent to godoc.org with the exception of /-/site.css. The site.css
from godoc.org is in g.go. Differences between godoc.org [site.css](https://godoc.org/-/site.css):
```css
.container { width: 100%; } 
/* .container { max-width: 728px; } */
```

Customize g.go site_css to your preferences.

#### License 

Use of this source code is governed by a BSD-2-Clause license that can be found
in the LICENSE file.

[![BSD-2-Clause License](osi_logo_100X133_90ppi_0.png)](https://opensource.org/)
