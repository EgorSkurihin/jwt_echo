//Set auth header for all requests
/*


button = document.getElementById("sayhello").onclick = function () {
    var xhr = new XMLHttpRequest();
    xhr.open('GET', '/restricted/sayhello', true);
    xhr.setRequestHeader('Authorization', 'Bearer ' + localStorage['token']);
    xhr.onload = function () {
        if (this.status == 401) {
            alert("You are not logged in")
        } else {
        console.log(this.responseText);
        }
    };
    xhr.send('');
};
*/
axios.interceptors.request.use(
    (config) => {
      const accessToken = localStorage['token'];
      config.headers.Authorization = `Bearer ${accessToken}`;
      return config;
    },
    (error) => Promise.reject(error)
);


document.getElementById("sayhello").onclick = function(){
    axios.get('http://localhost:1323/restricted/sayhello').then(resp => {
        //
        document.write(resp.data)
        });
}

document.getElementById("secretinfo").onclick = function(){
    axios.get('http://localhost:1323/restricted/secretinfo').then(resp => {
    //
    document.write(resp.data)
    });
}

//Auth on button 'login_button'
document.getElementById("login_button").onclick = function() {
    var req = new XMLHttpRequest();
    var name = document.getElementById("username").value;
    var password = document.getElementById("password").value;
    req.open('POST', 'http://localhost:1323/login', true);
    req.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
    req.onload = function () {
        if (this.status == 401) {
            alert("Wrong name or password!")
        }else{

            var token = JSON.parse(this.responseText)["token"];
            localStorage.setItem('token', token);
            window.location = 'http://localhost:1323/loggedin'
        }
    };
    req.send(`username=${name}&password=${password}`);
};