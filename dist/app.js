const store = new Vuex.Store({
    state: {
        lines: []
    },
    mutations: {
        addLine: function(state, line) {
            state.lines = [...state.lines, line]
        }
    }
})

const v = new Vue({
    el: "#app",
    store: store,
    computed: {
        year: function() {
            return new Date().getFullYear()
        },
        list: function() {
            return this.$store.state.lines
        }
    },
    template: '<div>\
        <section class="hero is-dark is-fullheight">\
            <div class="hero-body">\
                <div class="container">\
                    <h1 class="title">Roar</h1>\
                    <h2 class="subtitle">curl -i \'http://roar.io/api?cmd=./do_backend.sh&./do_front.sh\'</h2>\
                </div>\
            </div>\
        </section>\
        <section class="section">\
            <div class="container">\
                <h1 class="title">Output</h1>\
                <h2 class="subtitle">from all jobs</h2>\
                <hr />\
                <div class="content">\
                    <p>\
                        <template v-for="l in list">$ {{ l }}<br /></template>\
                    </p>\
                </div>\
            </div>\
        </section>\
        <section class="section has-text-centered">\
            BH {{ year }}\
        </section>\
    </div>',
})


let tryes = 0
function connectWebsocket() {
    if (tryes > 10) {
        alert("Não foi possível conectar ao websocket!");
        return
    }

    const wsURL = document.location.href.replace(/^http/, "ws")
    const ws = new WebSocket(wsURL + "ws")
    ws.onopen = function(e) {
        tryes = 0
    }
    ws.onclose = function() {
        setTimeout(connectWebsocket, (Math.random() * 2000 + 1000) | 0)
    }
    ws.onmessage = function(e) {
        const message = JSON.parse(e.data)
        store.commit("addLine", message.Payload)
    }
    tryes++
}

connectWebsocket();
