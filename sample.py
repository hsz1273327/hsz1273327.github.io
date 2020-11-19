"""牛顿法扩展都n次方根

:math:`\\sqrt[k] a` .的表达式为:

:math:`x_{n+1} = x_n - \\frac {x_n^k-a} {kx_n^{k-1}} = \\frac {k-1}{k} x_n + \\frac {a}{kx_n^{k-1}}`

"""


def sqrt_nt(n, m, round_=5):
    u"""牛顿法求n次方根.

    用法: :math:`sqrt\_nt(n,m,round) = \sqrt[m] n`


    Parameters:
        n (int,float,long): - 被开方的数
        m (int,float,long) : - 开多少次方
        round(int): - 精度

    Returns:
        int,float,long: 开方结果


    Raises:
        TypeError - 参数类型不对时报错
        ValueError - 被开方参数为复数时报错

    >>> sqrt_nt(2,2,round_=2)
    1.41

    """
    if all(map(lambda x: isinstance(x, (int, float, long,)), [n, m])):
        if n < 0:
            raise ValueError(u"必须是整数")
        elif n == 0:
            return 0
        elif n == 1:
            return 1
        else:
            deviation = 0.1**(round_ + 1)
            seed = n / m
            counter = 0
            max_count = 100000
            now_value = seed
            last_value = 0
            while abs(now_value - last_value) > deviation:
                if counter > max_count:
                    raise ValueError(u"在{count}次循环内未能得到精度为{round}的解".format(
                        count=max_count, round=round_))
                counter += 1
                last_value = now_value
                now_value = (1 - 1.0 / m) * last_value + \
                    n / (m * last_value**(m - 1))
            return round(now_value, round_)

    else:
        raise TypeError(u"必须是数")
