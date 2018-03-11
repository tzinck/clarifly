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
              <a v-on:click="createRoom" class="btn btn-default btn-lg">Create a Room</a>
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
        <img class="tree" src="./../assets/Treegreen.png">
        <img class="tree" src="./../assets/Treegreen.png">
        <img class="tree" src="./../assets/Treegreen.png">
        <img class="tree" src="./../assets/Treegreen.png">
        <img class="tree" src="./../assets/Treegreen.png">
        <img class="tree" src="./../assets/Treegreen.png">
        <img class="tree" src="./../assets/Treegreen.png">
        <img class="tree" src="./../assets/Treegreen.png">
        <img class="tree" src="./../assets/Treegreen.png">
        <img class="tree" src="./../assets/Treegreen.png">
        <img class="tree" src="./../assets/Treegreen.png">
        <img class="tree" src="./../assets/Treegreen.png">
        <img class="tree" src="./../assets/Treegreen.png">
        <img class="tree" src="./../assets/Treegreen.png">
        <img class="tree" src="./../assets/Treegreen.png">
        <img class="tree" src="./../assets/Treegreen.png">
      <div class="container about">
        <div class="row">
          <div class="col-lg-8 mx-auto">
            <h2>About Clarifly</h2>
            <p>Clarifly is a service that helps students ask questions in class without disturbing the lecture or feeling embarrassed.</p>
            <p>Teachers can create a room that students can submit questions to during class. The teachers can then answer the questions in class, or can view the questions after class and send an email to their class.</p>
            <p>Students can upvote eachother's questions to show that they also require clarification! </p>
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

        this.$http.post('http://localhost:8080/createRoom', {}).then(response => {


          // get body data
          //this.someData = response.body;
          this.someData = response.body.split(',')
          this.$store.commit('set_secret',this.someData[0]);
          this.$store.commit('set_room',this.someData[1]);
          this.room = this.someData[1];
          this.$store.commit('set_connected',true);
          //this.$router.push({ name: 'Join', params: { room: this.$store.state.room } });
          this.joinRoom();

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

        this.ws = new WebSocket("ws://localhost:8080/joinRoom");


        // On message: if room doesn't exist, close socket. 
        this.ws.onmessage = function(e) {

          if(String(e.data.trim()) === ("Room " + self.room + " does not exist.")){
            self.ws.close();
          }else{
            // Room exists, route to QuestionPage
            self.$store.commit('set_connected',true);

            var obj = JSON.parse(e.data);
            console.log(obj);
            if(obj.Code != undefined)
            {
              self.$store.commit('set_room', obj);
              console.log("meme");
              
            }
            else{
              console.log("meme2");
              self.$store.commit('set_question', obj);
              
            }
            self.$store.commit('set_ws', self.ws);
            self.$router.push({ name: 'Join', params: { room: self.$store.state.room.Code } });
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

input{
  width:203px;
}

#page-top{
  
  background-color: #F5D17B;
}

#about{
  color:white;
  font-weight:500;
  font-family: 'Cabin', 'Helvetica Neue', Helvetica, Arial, sans-serif;

}

.about{
  margin-top:20px;
}

.tree{
  width:50px;
}
</style>

