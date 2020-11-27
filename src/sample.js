/** @module sample */

/** The name of the module. */
export const name = 'sample'


/**
 * 测试函数
 * @param {string} color1 - 第一种颜色
 * @param {string} color2 - 第二种颜色
 * @return {string} 结果.
 */
function samplefunc (color1, color2) {
    return `${ color1 } annd ${ color2 }`
}

/** 二维点类型 */
class Point {
    /**
     * 创建一个二维点
     * @param {number} x - x轴的值.
     * @param {number} y - y轴的值.
     */
    constructor (x, y) {
        // ...
    }

    /**
     * 获取x轴的值
     * @return {number} x轴的值.
     */
    getX () {
        // ...
    }

    /**
     * 获取y轴的值.
     * @return {number} y轴的值.
     */
    getY () {
        // ...
    }

    /**
     * 从字符串解析生成点对象
     * @param {string} str - 待解析字符串
     * @return {Point} 点对象
     */
    static fromString (str) {
        // ...
    }
}

export {
    samplefunc,
    Point
}