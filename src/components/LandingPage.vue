<template>
    <div id="page-top">

    <!-- Intro Header -->
    <header class="masthead">
      <div class="intro-body">
        <div class="container">
          <div class="row">
            <div class="col-lg-8 mx-auto">
                <div class="logodiv">
                <img class="logo" src="../assets/logo.png" >
                </div>
               
              <p class="intro-text">Make your voice heard in class - without raising your hand</p>
              <a v-on:click="createRoom" class="btn btn-default btn-lg">Create a Room</a> {{someData}}
                <br>
            <a v-on:click="joinRoom" class="btn btn-join btn-lg">Join a Room</a>
            <br>
            <input maxlength="4" class="input-box" v-model="room" placeholder="room id">
            </div>
          </div>
        </div>
      </div>
    </header>

    <!-- About Section -->
    <section id="about" class="content-section text-center">
      <div class="container">
        <div class="row">
          <div class="col-lg-8 mx-auto">
            <h2>About Grayscale</h2>
            <p>Grayscale is a free Bootstrap theme created by Start Bootstrap. It can be yours right now, simply download the template on
              <a href="http://startbootstrap.com/template-overviews/grayscale/">the preview page</a>. The theme is open source, and you can use it for any purpose, personal or commercial.</p>
            <p>This theme features stock photos by
              <a href="http://gratisography.com/">Gratisography</a>
              along with a custom Google Maps skin courtesy of
              <a href="http://snazzymaps.com/">Snazzy Maps</a>.</p>
            <p>Grayscale includes full HTML, CSS, and custom JavaScript files along with SASS and LESS files for easy customization!</p>
          </div>
        </div>
      </div>
    </section>
    </div>
</template>

<script>
import{mapState, mapMutations} from 'vuex';

export default {
  name: 'LandingPage',
  data () {
    return {
      ws: null,
      msg: 'Welcome to Your Vue.js App',
      someData: "",
      room: ""
    }
  },
  // computed: {
  //   ...mapState(['connected'])
  // },
  methods: {
      createRoom() {
          console.log("mememe");
        // GET /someUrl
        this.$http.post('http://9081f8c8.ngrok.io/createRoom', {}).then(response => {

          // get body data
          this.someData = response.body;
          this.$store.commit('set_secret',this.someData[0]);
          this.$store.commit('set_room',this.someData[1]);
          
          this.$router.push({ name: 'Join', params: { room: this.$store.state.room } });

        }, response => {
          console.log(response);
          // error callback
        });
      },

      joinRoom() {
          console.log(this.room);
        // GET /someUrl
        //this.$http.post('http://889a3db6.ngrok.io/joinRoom', {room_id: this.room}).then(response => {
          var self = this;
          console.log(this.$store.state.ws.readyState);

        // Check to see if the ws already exists and that it's not closed
        if (this.$store.state.ws && this.$store.state.ws.readyState == '1') {
          // Exists, route to it
          console.log("here");
          self.$router.push({ name: 'Join', params: { room: self.$store.state.room } });
        }else{
        // Doesn't exist, try to make a new one

        // Reset vars
        self.$store.commit('set_connected',false);
        self.$store.commit('set_room', '');
        self.$store.commit('set_ws', '');

        // Open websocket
        this.ws = new WebSocket("ws://9081f8c8.ngrok.io/joinRoom");

        // On message: if room doesn't exist, close socket. 
        this.ws.onmessage = function(e) {

          if(String(e.data.trim()) === ("Room " + self.room + " does not exist.")){
            self.ws.close();
          }else{
            // Room exists, route to QuestionPage
            self.$store.commit('set_connected',true);
            self.$store.commit('set_room', JSON.parse(e.data));
            self.$store.commit('set_ws', self.ws);
            self.$router.push({ name: 'Join', params: { room: self.$store.state.room } });
          }
        };
          // On open: send our room code
          this.ws.onopen = function() {
            self.ws.send(self.room);
          };
        }

        

          // get body data
          // this.someData = response.body;
          // this.set_room(this.room);
      }
    }
}
</script>
<style scoped>

.intro-text{
  font-weight:500;
}

body{
  background-color:#90D0E5;
}
</style>

