var app = new Vue({
    el: '#main',
    delimiters: ["[[", "]]"],
    data: function() {
        return {}
    },
    methods: {
        newshow: function() {
            console.log("new show")
            window.location.replace("/show/new")
        }
    }
})