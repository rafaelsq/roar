const store = new Vuex.Store({
    state: {
        actives: [],
        lines: []
    },
    mutations: {
        addLine: function(state, line) {
            state.lines = [...state.lines, line]
        },
        addActives: function(state, active) {
            state.actives = [...state.actives, active]
        }
    }
})

const cmd = {
    template: '\
       <section class="section">\
           <div class="container">\
               <h1 class="title">Output</h1>\
               <h2 class="subtitle">from all jobs</h2>\
               <hr />\
               <div class="content">\
                   <div class="columns">\
                        <div class="column" v-for="c in active.Commands">\
                            <h3>{{ c }}</h3>\
                           <div v-for="l in filterByCommand(c)" :key="l.Payload.Id" :class="{notification: l.Type > 0, \'is-danger\': l.Type == 1, \'is-success\': l.Type ==2}">{{ l.Payload.Text }}</div>\
                        </div>\
                   </div>\
                   <div v-for="l in listCommon" class="notification" :class="{\'is-danger\': l.Type == 1, \'is-success\': l.Type ==2}">{{ l.Payload.Text }}</div>\
               </div>\
           </div>\
       </section>',
    data: () => ({
        Id: 0
    }),
    mounted: function() {
        this.Id = parseInt(this.$route.params.Id, 10)
    },
    watch: {
        '$route' (to, from) {
            this.Id = parseInt(to.params.Id, 10)
        }
    },
    methods: {
        filterByCommand: function(cmd) {
            return this.items.filter((o) => o.Payload.Command === cmd)
        }
    },
    computed: {
        listCommon: function() {
            return this.items.filter((o) => !o.Payload.Command)
        },
        items: function() {
            return this.$store.state.lines.filter((o) => o.Payload.Id === this.Id)
        },
        active: function() {
            if (!this.Id) {
                return {Commands: []}
            }

            const matchs = this.$store.state.actives.filter((a) => a.Id === this.Id)
            if (matchs.length) {
                return matchs[0]
            }

            return {Commands: []}
        }
    }
}

new Vue({
    store: store,
    data: () => ({
        isHome: false,
    }),
    computed: {
        year: function() {
            return new Date().getFullYear()
        },
        actives: function() {
            return this.$store.state.actives.map((a) => ({Id: a.Id, Path: "/cmd/" + a.Id}))
        }
    },
    mounted: function() {
        this.isHome = this.$router.history.current.path === '/'
    },
    watch: {
        '$route': function(to) {
            this.isHome = to.path === '/'
        }
    },
    router: new VueRouter({
        mode: 'history',
        linkActiveClass: 'is-active',
        routes: [
            {path: '/'}, 
            {path: '/cmd/:Id', component: cmd}, 
        ]
    }),
    template: '<div>\
        <section class="hero is-dark" v-bind:class="{\'is-fullheight\': isHome}">\
            <div v-if="isHome" class="hero-body">\
                <div class="container">\
                    <h1 class="title">Roar</h1>\
                    <h2 class="subtitle">curl -i \'http://roar.io/api?cmd=./do_backend.sh&./do_front.sh\'</h2>\
                </div>\
            </div>\
            <div class="hero-foot">\
                <nav class="nav has-shadow">\
                    <div class="container">\
                        <div class="nav-center">\
                            <router-link class="nav-item is-tab" exact to="/">home</router-link>\
                            <router-link v-for="a in actives" :key="a.Id" class="nav-item is-tab" :to="a.Path">#{{ a.Id }}</router-link>\
                        </div>\
                    </div>\
                </nav>\
            </div>\
        </section>\
        <router-view></router-view>\
        <section class="section has-text-centered">\
            BH {{ year }}\
        </section>\
    </div>'
}).$mount("#app")

let tryes = 0
function connectWebsocket() {
    if (tryes > 10) {
        alert("Não foi possível conectar ao websocket!");
        return
    }

    const l = document.location
    const ws = new WebSocket(l.protocol.replace(/^http/, "ws") + "//" + l.host + '/ws')
    ws.onopen = function(e) {
        tryes = 0
    }
    ws.onclose = function() {
        setTimeout(connectWebsocket, (Math.random() * 2000 + 1000) | 0)
    }
    ws.onmessage = function(e) {
        const msg = JSON.parse(e.data)
        if (msg.Type == 3) {
            store.commit("addActives", msg.Payload)
        } else {
            store.commit("addLine", msg)
        }
    }
    tryes++
}

connectWebsocket()
