import sphinx_rtd_theme
from recommonmark.transform import AutoStructify
from recommonmark.parser import CommonMarkParser
project = 'sample'
copyright = '2020, hsz'
author = 'hsz'

# 确定项目版本
version = '0.0.0'
release = '0.0.0'

# 使用的插件
extensions = [
    'sphinx.ext.todo',
    'sphinx.ext.mathjax',
    'sphinx.ext.viewcode',
]

# 入口文件
master_doc = 'index'

# 指定模板所在文件夹位置
templates_path = ['_templates']

# 文档语言
language = 'zh_CN'

# 不进行编译的文件/文件夹
exclude_patterns = ['_build', 'Thumbs.db', '.DS_Store']

# 设置不同后缀的文件使用不同解析器(这个需要后加)
source_suffix = {
    '.rst': 'restructuredtext'
}


# 指定编译成html时使用的主题
html_theme = 'alabaster'

# 指定编译成html时使用的静态文件所在位置
html_static_path = ['_static']


# todo插件的设置
todo_include_todos = True

# 主题
extensions.append('sphinx_rtd_theme')
html_theme = "sphinx_rtd_theme"
html_theme_options = {
    'logo_only': True,
    'navigation_depth': 5,
}


# 使用插件支持markdowm
extensions.append('recommonmark')

# 针对`.md`为后缀的文件做markdown渲染
source_suffix[".md"] = 'markdown'

# 设置markdown渲染器的自定义项


def setup(app):
    github_doc_root = 'https://localhost:5000'
    app.add_config_value('recommonmark_config', {
        # 'url_resolver': lambda url: github_doc_root + url,
        "enable_auto_toc_tree": True,
        "auto_toc_tree_section": "目录",
        'auto_toc_maxdepth': 2,
        "enable_math": True,
        'enable_eval_rst': True
    }, True)

    app.add_transform(AutoStructify)


# C 语言支持
extensions.append('breathe')
extensions.append('exhale')
# Setup the breathe extension
breathe_projects = {
    "Sample": "./doxyoutput/xml"
}

breathe_default_project = "Sample"

# Setup the exhale extension
exhale_args = {
    # These arguments are required
    "containmentFolder": "./api",
    "rootFileName": "index.rst", # "library_root.rst",
    "rootFileTitle": "Library API",
    "doxygenStripFromPath": "..",
    # Suggested optional arguments
    "createTreeView": True,
    # TIP: if using the sphinx-bootstrap-theme, you need
    # "treeViewIsBootstrap": True,
    "exhaleExecutesDoxygen": True,
    "exhaleDoxygenStdin": "INPUT = ../include"
}

# Tell sphinx what the primary language being documented is.
primary_domain = 'cpp'

# Tell sphinx what the pygments highlight language should be.
highlight_language = 'cpp'
