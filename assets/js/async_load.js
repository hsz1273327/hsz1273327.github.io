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