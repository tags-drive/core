new Vue({
    el: "#login-form",
    data: {
        login: "",
        password: "",
        isError: false,
        errorMsg: ""
    },
    methods: {
        fail: function(error) {
            setTimeout(() => {
                this.errorMsg = error;
                this.isError = true;
                this.password = "";
            }, 700);
        },
        auth: function() {
            this.isError = false;

            // 11 times
            var hash = sha256(this.password);
            for (var i = 0; i < 10; i++) {
                hash = sha256(hash);
            }

            const params = new URLSearchParams();
            params.append("login", this.login);
            params.append("password", hash);

            fetch("/login", {
                method: "POST",
                body: params,
                credentials: "same-origin" // for set cookie
            })
                .then(data => {
                    // Valid login and password
                    if (data.status == 200) {
                        window.location.href = "/";
                    }

                    data.text().then(msg => this.fail(msg));
                })
                .catch(err => console.log(err));
        }
    }
});
