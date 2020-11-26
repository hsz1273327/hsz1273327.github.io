# 使用Sphinx写项目文档的演示代码

本例子演示使用markdown写多语言支持的文档


使用方法:

> 多语言支持

+ 将文档内容转化为`pot`文件,放入`document/_build/gettext`文件夹下

    ```bash
    sphinx-build -b gettext document document/_build/gettext
    ```

+ 使用`sphinx-intl`工具将`pot`文件都转化成`po`文件用于翻译

    ```bash
    sphinx-intl update -p document/_build/gettext -d document/locale -l 语言1 -l 语言2 ...
    ```

    这一步的`-d document/locale`和上面配置的`locale_dirs`对应

> 编译成html

+ 编译原始语言的html执行
    + windows`$env:CURRENT_LANGUAGE="zh_CN"; sphinx-build -b html document docs`
    + linux/macos`export CURRENT_LANGUAGE=zh_CN && sphinx-build -b html document docs`

+ 编译成不同语言的html执行`sphinx-build document docs`后在docs文件夹下查看效果
    + windows`$env:CURRENT_LANGUAGE="fr"; sphinx-build -D language='fr' -b html document docs/fr`
    + linux/macos`export CURRENT_LANGUAGE=fr && sphinx-build -D language='fr' -b html document docs/fr`