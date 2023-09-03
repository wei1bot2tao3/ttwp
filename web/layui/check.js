
 const tokenc = localStorage.getItem('token');
var username1=localStorage.getItem('username')

if (!tokenc||!username1) {
    // 如果不存在token，则需要用户重新登录
    setTimeout(function() {
        layer.msg('请先登录', {icon: 2});
    });
   window.location.href = '/';

}

// 向服务器发送请求，验证token是否有效
fetch('/check', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
    },
    body: new URLSearchParams({
        'username': username1,
        'token': tokenc
    })



})
    .then(response => response.json())
    .then(data => {
        if (data.code === 1) {
            // 如果token有效，则执行您需要的操作
            layer.msg('token有效验证登录成功');
        } else {
            // 如果token无效，则需要用户重新登录
            setTimeout(function() {
                alert("登录已过期，请重新登录")
                window.location.href = '/';
            }, 10000);

        }
    })
    .catch(error => {
        console.error(error);
        alert('发生错误，请稍后再试');
    });



