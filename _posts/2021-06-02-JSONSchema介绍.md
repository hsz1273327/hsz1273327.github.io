---
layout: post
title: "JSONSchema介绍"
date: 2021-06-02
author: "Hsz"
category: introduce
tags:
    - WebTech
    - DataTech
    - DataAnnotate 
    - DataValidate 
header-img: "img/home-bg-o.jpg"
update: 2021-06-02
---
# JSONSchema介绍

[JSONSchema](http://json-schema.org/)是用来定义JSON数据约束的一个标准.根据这个约定模式,交换数据的双方可以理解JSON数据的要求和约束,也可以据此对数据进行验证以保证数据交换的正确性.
当然由于许多JSONSchema实现都附带可以校验用于生成JSON的对应结构体或者哈希表结构等,它事实上也可以用于做结构化数据校验.

目前最新的JSONSchema版本是[Draft 8](http://json-schema.org/specification-links.html#2019-09-formerly-known-as-draft-8),但多数编程语言的实现目前只支持到[draft-4](http://json-schema.org/specification-links.html#draft-4).第二多的则是支持到[draft-7](http://json-schema.org/specification-links.html#draft-7)

本文将以draft-4为基准介绍JSONSchema.并顺便在碰到draft-7新增或者变动的内容时提一下

## 一个典型的JSONSchema例子

```json
{
  "$id": "https://example.com/address.schema.json",
  "$schema": "http://json-schema.org/draft-04/schema#",
  "description": "An address similar to http://microformats.org/wiki/h-card",
  "type": "object",
  "properties": {
    "post-office-box": {
      "type": "string"
    },
    "extended-address": {
      "type": "string"
    },
    "street-address": {
      "type": "string"
    },
    "locality": {
      "type": "string"
    },
    "region": {
      "type": "string"
    },
    "postal-code": {
      "type": "string"
    },
    "country-name": {
      "type": "string"
    }
  },
  "required": [ "locality", "region", "country-name" ],
  "dependencies": {
    "post-office-box": [ "street-address" ],
    "extended-address": [ "street-address" ]
  }
}
```

JSON Schema本身也是JSON,它是可以[自举](https://baike.baidu.com/item/%E8%87%AA%E4%B8%BE/3063932?fr=aladdin)的,draft-4的自描述文件在<http://json-schema.org/draft-04/schema#>,我们可以对着下面的介绍看这个自描述文件.

## 概念

JSONSchema规范中定义的一些概念用于将这个规范描述清楚,他们分别是:

+ 实例(instance)

    表示应用于JSONSchema的JSON文档,实例中的字段我们称之为实例的关键字(instance keyword)

+ 模式(Schema)

    用于描述实例或实例关键字应该满足的条件集合.JSON Schema本身就是一个模式,其中对每个关键字的描述都是一个模式.模式是可以单独定义多处复用的,这个后面会有介绍.

+ 验证关键字(Validation Keyword)

    模式可用于要求给定的实例或实例的关键字满足一定数量的条件,而用于断言这些条件的关键字则被称为验证关键字.验证关键字根据其适用范围分为如下几类:

    + 适用于任何实例类型的验证关键字
    + 数值实例的验证关键字
    + 字符串实例的验证关键字
    + 数组实例的验证关键字
    + 对象实例的验证关键字
    + 有条件的使用子模式的关键字
    + 根据布尔逻辑使用子模式的关键字
    + 用于验证格式语意的关键字

    基本上定义一个JSONSchema就是使用验证关键字描述实例的过程.

+ 注释关键字(Schema Annotations)

    JSONSchema中还可以使用其他一些用于描述JSONSchema本身的注释关键字,通常我们希望让JSONSchema更具可读性就会添加这些关键字

+ 模式关键字(Schema Keyword)

    一个JSONSchema必须是一个对象或者布尔值,声明这个模式本身的关键字就是模式关键字

## 常见的注释关键字

+ `title`取值必须是字符串类型,表示模式的名称,应该尽量简短.

+ `description`取值必须是字符串类型用于描述模式的用途简介.

+ `default`指定项的默认值.注意这个默认值的含义是建议使用指定值作为默认值,一般的验证器实现都不会取实现这个默认值自动补正的功能.

+ `examples`取值必须是一个数组,提供一组满足模式的实例用于向读者解释模式的作用和效果.

+ `readOnly`/`writeOnly`(draft-7)取值必须是布尔值.如果`readOnly`的值为`true`则表明实例的值仅由拥有者管理,并且应用程序尝试修改此属性的值预计将被拥有者忽略或拒绝;如果`writeOnly`的值为`true`则表示从拥有权限中检索实例时该值将永远不存在.
    例如`readOnly`将用于标记数据库生成的序列号为只读.而`writeOnly`将用于标记密码输入字段.

+ `$comment`(draft-7)取值必须为字符串,表示纯粹的注释

## 常见的模式关键字

+ `id`/`$id`(draft-7)用于描述这个JSON Schema的唯一标识符,其唯一性是通过网络环境中一个uri来确定的.
+ `$schema`用于描述这个JSON Schema使用的是哪一版协议,像draft-7就应该填入`"http://json-schema.org/draft-07/schema#"`

## 布尔值作为模式

布尔型数值本身可以作为模式使用,其含义为:

+ `true`等价于`{}`,既无论什么实例都可以通过验证
+ `false`等价于`{"not":{}}`,既无论什么实例都不通过验证

## 验证关键字

### 通用的验证关键字

+ `type`这个关键字的值必须是字符串或数组.当且仅当实例位于为该关键字列出的任何集合中时实例才验证成功.

    + 如果它是一个字符串值必须是六种基本类型(`null`,`boolean`,`object`,`array`,`number`,`string`)中的一种或者是表示整数的`integer`

        例如,如果要验证一个字符串类类型的实例，可以使用如下的模式:

        ```json
        { "type": "string" }
        ```

    + 如果它是一个数组,那么数组的元素必须是字符串,必须是惟一的,且必须上上面7种描述字符串之一
        例如,如果要验证一个字符串类型或者数值类型或者数组类型的实例，可以使用如下模式:

        ```json
        { "type": ["string", "number", "array"] }
        ```

+ `enum` 这个关键字的值必须是一个数组,这个数组应该至少有一个元素且数组中的元素应该是惟一.它代表实例关键字的取值范围必须落在这个数据所包含的元素中.

    如果实例的值等于该关键字数组值中的某一个元素则实例验证成功.数组中的元素可以是任何值包括`null`.这个验证关键之一般用于缩小实例关键字的取值范围
    例如要验证一个实例的实例关键字`country`是否在`"China","UK","US"`这三个值的范围内就可以使用如下的模式:

    ```json
    {
        "country":{
            "type": "string",
            "enum": ["China", "UK", "US"]
        }
    }
    ```

+ `const`(draft-7)这个关键字的值可以是任何类型包括`null`.它表示这个实例关键字必须是一个指定的值.

    如果一个实例的值等于这个关键字的值则验证成功.
    例如要验证一个实例的实例关键字`country`是否是`China`,可以使用如下的模式:

    ```json
    {
        "country":{
            "const": "China"
        }
    }
    ```

### number类型的专用验证关键字

+ `multipleOf`取值必须是一个大于0的数,这个验证关键字的含义是可以被这个设定值整除

+ `maximum`取值必须是一个数,表示实例的关键字值必须小于等于设定值

+ `exclusiveMaximum`取值必须是一个数,表示实例的关键字值必须小于设定值

+ `minimum`取值必须是一个数,表示实例的关键字值必须大于等于设定值

+ `exclusiveMinimum`取值必须是一个数,表示实例的关键字值必须大于设定值

### string类型的专用验证关键字

+ `maxLength`取值必须是一个非负整数,这个验证关键字的含义是字符串长度需要小于等于这个设定值

+ `minLength`取值必须是一个非负整数,这个验证关键字的含义是字符串长度需要大于等于这个设定值,如果这个关键字的值为0则相当于忽略这个关键字

+ `pattern`取值必须是字符串.根据ECMA 262正则表达式方言,这个字符串应该是一个有效的正则表达式.这个验证关键字的含义是取值必须满足正则表达式的匹配要求.

+ `format`取值必须是字符串,这个可以理解为是一些常用特殊`pattern`的快捷方式.目前规定的有:
    + 日期和时间
        + `date`
        + `time`
        + `datetime`
    + 电子邮件地址
        + `email`
        + `idn-email`
    + 主机名
        + `hostname`
        + `idn-hostname`(draft-7)
    + IP地址
        + `ipv4`
        + `ipv6`
    + 资源标识符
        + `uri`
        + `uri-reference`(draft-7)
        + `iri`
        + `iri-reference`(draft-7)
    + URI模版
        + `uri-template`(draft-7)
    + JSON指针
        + `json-pointer`
        + `relative-json-pointer`(draft-7)
    + 正则表达式
        + `regex`

    但需要注意,对`format`的支持并不是所有的验证其实现都有,有的即便有也需要安装额外依赖.很多时候需要看下使用的验证器.

### array类型的专用验证关键字

+ `items`取值必须是一个有效的模式或一个有效的模式数组.这个验证关键字用于描述array中的元素.

    + 如果这个关键字的值是一个模式,那么只有数组中的所有元素都对该模式成功验证才算验证成功.
        例如要验证数组实例的每个元素是否都是长度大于等于3且小于等于5的字符串实例,可以使用如下的模式:

        ```json
        {
            "type": "array",
            "items": {
                "type": "string",
                "minLength": 3,
                "maxLength": 5
            }
        }
        ```

    + 如果这个关键字是模式数组,那么如果实例的每个元素在相同的位置(如果有的话)对模式成功验证才算验证成功.

        例如，要验证数组实例中至少有3个元素且前3个元素是否依次为字符串,数值和布尔类型,可以使用如下的模式:

        ```json
        {
            "type": "array",
            "items": [
                {"type": "string"},
                {"type": "number"},
                {"type": "boolean"}
            ]
        }
        ```

+ `additionalItems`取值必须是一个有效的模式
    这个验证关键字需要配合`items`使用,且只有在`items`取值为模式数组时才可以声明,否则无效.它的含义是在`items`规定的模式数组之外的元素可以使用这个字段定义的模式来验证.

+ `maxItems`取值必须是一个非负整数,含义是数组中的元素个数需要小于等于这个设定值

+ `minItems`取值必须是一个非负整数,含义是数组中的元素个数需要大于等于这个设定值

+ `uniqueItems`取值必须是一个布尔值,含义为数组中的每个元素是否唯一

+ `contains`(draft-7)取值必须是一个有效的JSON Schema.它的含义是数组中的元素必须有至少一个可以匹配设定的模式.

    例如如果要验证一个数组实例的个数大于等于3且小于等5并且数组中的每个元素都是唯一的并且至少有一个元素是数字0,则可以使用如下的模式:

    ```json
    {
        "type": "array",
        "minItems": 3,
        "maxItems": 5,
        "uniqueItems": True,
        "contains": {
            "type": "number",
            "const": 0
        }
    }
    ```

### object类型的专用验证关键字

+ `maxProperties`取值必须是一个非负整数,含义是对象的属性个数要小于等于这个设定值.

+ `minProperties`取值必须是一个非负整数,含义是对象的属性个数要大于等于这个设定值.

+ `required`取值必须是一个字符串为元素的数组,且元素必须唯一.含义是对象的必须包含设定值中的全部属性.

+ `properties`取值必须是一个对象,对象的每个属性为要用来做匹配的实例中的属性名,值为对应的模式.其含义是实例的对应属性需要满足对应的模式

    例如要验证如果一个对象实例中包含字符串型的`name`和数值型的`age`属性可以使用如下的模式:

    ```json
    {
    "type": "object",
    "properties": {
        "name": {
            "type": "string"
        },
        "age": {
            "type": "number"
        }
    }
    }
    ```

    注意这个例子的含义并不是只能有`name`和`age`两个属性,如果匹配的对象中还有比如`country`它同样可以通过.这个验证关键字只会验证其中定义了的字段.

+ `patternProperties`取值必须是一个对象,该对象的每个属性名都应该是一个有效的正则表达式,该对象的每个属性的值必须是一个有效的模式.这个验证关键字是上面`properties`的补充和加强,它可以通过对实例对象属性的正则匹配找到要验证模式的属性进行验证.通常用于对属性做批量校验.
    例如如果要验证对象实例中所有属性名以下划线开头的属性的属性值的长度是否都大于等于3并且所有属性名以下划线结束的属性的属性值是否都是10的整数倍则可以使用如下的模式:

    ```json
    {
        "type": "object",
        "properties": {
            "_name": {
                "type": "string"
            },
            "age_": {
                "type": "number"
            },
            "salary_": {
                "type": "number"
            },
            "_department": {
                "type": "string"
            }
        },
        "patternProperties": {
            "^_": {
                "minLength": 3
            },
            "_$": {
                "multipleOf": 10
            }
        }
    }
    ```

+ `additionalProperties`取值必须是一个模式.含义是实例中在`properties`关键字中定义的属性以及`patternProperties`关键字能匹配到的属性外的其他属性应该满足什么样的模式.不过最常见的用法是直接用`false`作为模式禁止`properties`关键字和`patternProperties`关键字定义之外的属性存在.

+ `dependencies`取值必须是一个对象,其属性为需要指定依赖的实例属性,属性的值为模式或者元素为字符串型的array.它用于描述属性的依赖关系.

    + 如果属性的值为元素为字符串型的array,那么这个字符串必须唯一且都是实例中包含的属性
        例如我们要用模式表示一个信用卡用户,当用户有信用卡卡号时必须同时存在账单地址,那么模式可以写成:

        ```json
        {
            "type": "object",
            "properties": {
                "name": { "type": "string" },
                "credit_card": { "type": "number" },
                "billing_address": { "type": "string" }
            },
            "required": ["name"],
            "dependencies": {
                "credit_card": ["billing_address"]
            }
        }
        ```

    + 如果属性的值为模式,则这个模式一定是描述对象的模式,含义是实例中需要有对应的字段且对应字段满足特定模式.

        还是上面的例子,我们也可以写成:

        ```json
        {
            "type": "object",
            "properties": {
                "name": { "type": "string" },
                "credit_card": { "type": "number" }
            },
            "required": ["name"],
            "dependencies": {
                "credit_card": {
                    "properties": {
                        "billing_address": { "type": "string" }
                    },
                    "required": ["billing_address"]
                }
            }
        }
        ```

+ `propertyNames`(draft-7)取值必须是一个有效的模式,其含义是实例的对象中属性名必须都满足这个设定的模式.由于属性名本来就是string,因此后面的模式可以不额外声明`type`.

## 模式的条件匹配

我们可能会有这样的需求--如果`locale`为`CN`则我们的`language`可以为`ZH`或者`EN`,否则`language`只能是`EN`,这里面就涉及到了条件判断,我们可以使用验证关键字组合`if`,`then`,`else`,这三个关键字的值都是模式,并且一定要配合使用.上面的例子我们可以使用模式:

```json
{
  "type": "object",
  "properties": {
    "locale": {
      "type": "string"
    }
  },
  "if": {
    "properties": { "locale": { "const": "CN" } }
  },
  "then": {
    "properties": { "language": { "enum": ["ZH","EN"] } }
  },
  "else": {
    "properties": { "language": { "const": "EN" } }
  }
}
```

## 布尔逻辑验证模式

有的时候光用一个模式来描述实例会不够,这种时候我们可以通过使用布尔逻辑验证关键字来规定验证工具的行为.

+ `allOf`取值必须为一个模式列表,它的含义是实例需要满足所有设置的模式.这个一般用的不多.

+ `anyOf`取值必须为一个模式列表,它的含义是实例需要满足至少一个设置的模式.

+ `oneOf`取值必须为一个模式列表,它的含义是实例需要满足其中一个设置的模式.通常`oneOf`中的模式应该是互斥的.

+ `not`取值必须是一个模式,其含义是如果实例不满足设置的模式则通过验证.

一个常见的用法是需要描述的实例可以有多种情况,比如我们允许`language`字段是`CN`,`EN`,`FR`之中的一个或几个,那么可以这样描述:

```json
{
    "type": "object",
    "properties": {
        "language": {
            "oneOf":[{
                "type":"string",
                "enum":["CN","EN","FR"]
            },{
                "type":"array",
                "items":{
                    "type":"string",
                    "enum":["CN","EN","FR"]
                }
            }]
        }
    }
}
```

## 模式复用

模式的复用可以有两种方式,一种是使用内部定义的子模式,另一种是使用外部定义的模式.他们的调用方式都是使用模式关键字`$ref`,其取值为一个相对或者绝对URI`#`符号后面的部分是JSON指针.

### 子模式

在JSON Schema中我们可以将共用的部分抽取成子模式,这样就可以在模块内部直接复用了.定义子模式的方式是使用模式关键字`definitions`/`$defs`(draft-7),其取值为一个对象,对象的属性为子模式的名字,值为子模式.
比如上面的例子我们就可以利用子模式抽出共用的部分:

```json
{
    "type": "object",
    "$defs":{
        "language_string":{
            "type":"string",
            "enum":["CN","EN","FR"]
        }
    },
    "properties": {
        "language": {
            "oneOf":[{
                "$ref":"#/$defs/language_string"
            },{
                "type":"array",
                "items":{
                    "$ref":"#/$defs/language_string"
                }
            }]
        }
    }
}
```

### 使用外部模式

JSON基本算是为网络而生,JSON Schema也基本上如此,因此很多时候我们会将模式定义好后直接挂在网上,当需要使用时直接在`$ref`后面写上要用的地址即可,比如官网的例子就是,我们借用官网的例子[geographical-location.schema.json](http://json-schema.org/learn/examples/geographical-location.schema.json)来演示使用外部模式.

我们要定义一个用户信息,其定义如下:

```json
{
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "properties": {
        "name": {
            "type":"string"
        },
        "age":{
            "type":"integer"
        },
        "location":{
            "$ref":"http://json-schema.org/learn/examples/geographical-location.schema.json"
        }
    }
}
```

然后我们随便写个实例验证下:

```json
{
  "name": "asf",
  "age": 12,
  "location": {
    "latitude": 32.2,
    "longitude": 43
  }
}
```

## 常见编程语言的校验器

| 编程语言   | 模块名                                                                     | 注意点                                                                                                     |
| ---------- | -------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- |
| python     | [jsonschema](https://github.com/Julian/jsonschema)                         | 默认没有format校验,但可以配置添加;访问外部模式使用`requests`,没有异步访问的实现,因此可能需要自己写猴子补丁 |
| javascript | [ajv](https://github.com/epoberezkin/ajv)                                  | ---                                                                                                        |
| golang     | [github.com/xeipuuv/gojsonschema](https://github.com/xeipuuv/gojsonschema) | ---                                                                                                        |
| C++        | [jsoncons](https://github.com/danielaparker/jsoncons)                      | ---                                                                                                        |

如果只是为了验证语法,推荐直接使用[网页版的校验器](https://www.jsonschemavalidator.net/)
