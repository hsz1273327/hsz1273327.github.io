// 导航栏
// Drop Bootstarp low-performance Navbar
// Use customize navbar with high-quality material design animation
// in high-perf jank-free CSS3 implementation
var $body = document.body;
var $toggle = document.querySelector('.navbar-toggle');
var $navbar = document.querySelector('#huxblog_navbar');
var $collapse = document.querySelector('.navbar-collapse');

$toggle.addEventListener('click', handleMagic)

function handleMagic(e) {
    if ($navbar.className.indexOf('in') > 0) {
        // CLOSE
        $navbar.className = " ";
        // wait until animation end.
        setTimeout(function () {
            // prevent frequently toggle
            if ($navbar.className.indexOf('in') < 0) {
                $collapse.style.height = "0px"
            }
        }, 400)
    } else {
        // OPEN
        $collapse.style.height = "auto"
        $navbar.className += " in";
    }
}
// resize header to fullscreen keynotes
var $header = document.getElementsByTagName("header")[0];
var $header_container = document.getElementsByClassName("container")[0];

function resize() {
    /*
     * leave 85px to both
     * - told/imply users that there has more content below
     * - let user can scroll in mobile device, seeing the keynote-view is unscrollable
     */
    // $header.style.height = (window.innerHeight - 85) + 'px';
    $header.style.height = '40vh';
    $header_container.height = '40vh'
}
document.addEventListener('DOMContentLoaded', function () {
    resize();
})
window.addEventListener('load', function () {
    resize();
})
window.addEventListener('resize', function () {
    resize();
})
resize();

// async load function
function async_load(u, c) {
    var d = document,
        t = 'script',
        o = d.createElement(t),
        s = d.getElementsByTagName(t)[0];
    o.src = u;
    if (c) {
        o.addEventListener('load', function (e) {
            c(null, e);
        }, false);
    }
    s.parentNode.insertBefore(o, s);
}

// only load tagcloud.js in tag.html
if ($('#tag_cloud').length !== 0) {
    async_load("/assets/js/jquery.tagcloud.js", function () {
        $.fn.tagcloud.defaults = {
            //size: {start: 1, end: 1, unit: 'em'},
            color: {
                start: '#bbbbee',
                end: '#0085a1'
            },
        };
        $('#tag_cloud a').tagcloud();
    })
}
// fastClick.js 解决移动端延迟
async_load("https://cdn.bootcss.com/fastclick/1.0.6/fastclick.min.js", function () {
    var $nav = document.querySelector("nav");
    if ($nav) FastClick.attach($nav);
})
