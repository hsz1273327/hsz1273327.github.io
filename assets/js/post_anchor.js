async_load("http://cdn.bootcss.com/anchor-js/1.1.1/anchor.min.js",function(){
    anchors.options = {
        visible: 'always',
        placement: 'left',
        icon: '¶'
    };
    anchors.add().remove('.intro-header h1').remove('.subheading').remove('.sidebar-container h5');
})