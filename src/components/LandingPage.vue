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
              <v-btn @click.native="createRoom" class="btn btn-default btn-lg">Create a Room</v-btn>
                <br>
            <a href="" class="btn btn-join btn-lg">Join a Room</a>
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
export default {
  name: 'HelloWorld',
  data () {
    return {
      msg: 'Welcome to Your Vue.js App'
    }
  },
  methods: {
      createRoom: function() {
        this.ajaxRequest = true;
        this.$http.post('http://28ffb55a.ngrok.io/createRoom', {},
         function (data, status, request) {
                this.postResults = data;
                this.ajaxRequest = false;
            });
      }
    }
}
</script>
<style scoped>
body{
  background-color:#90D0E5;
}
</style>

