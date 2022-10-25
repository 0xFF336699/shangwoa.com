```aidl
{
  "map": {
    "baobeidaola-android-web-shell": {
	  "prefix": "baobeidaola-android-web-shell",
	  "filePaths": [
		{
		  "directory": "/app/baobeidaola-android-web-shell/"
		}
	  ],
      "hosts": [
		"shell.web.fanli.app.baobeidaola.com"
      ],
      "redirectMap": {
        "/www/": "/www/index.html"
      }
	}
  }
}


必须要配置hosts

"directory": "/app/baobeidaola-android-web-shell/"
"/www/": "/www/index.html"
整个路径可能是
/app docker虚拟路径
/baobeidaola-android-web-shell/ docker容器app目录下的目录

http://x.com/www/
实际访问的是docker内路径
/app/baobeidaola-android-web-shell/www/index.html
```
