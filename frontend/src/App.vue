<script setup>
import { RouterLink, RouterView } from "vue-router";
</script>

<template>
  <header>
    <div class="container">
      <div class="row">
        <div class="one-half column">
          <RouterLink to="/">
            <img alt="MaxPanel home page" class="logo" src="./assets/logo.svg" width="300" height="100" />
          </RouterLink>
        </div>
        <div class="one-half column menu">
          <nav>
            <RouterLink to="/signup/" class="button button-primary signup">Sign Up</RouterLink>
            <RouterLink to="/login/" class="button button-primary">Login</RouterLink>
          </nav>
        </div>
      </div>
      
    </div>
  </header>

  <RouterView />
</template>

<style>
@import url('https://fonts.googleapis.com/css?family=Raleway:400,300,600');

body {
  background-color: #64634d;
}
.menu {
  text-align: right;
  margin-top: 3rem;
}
a.button.button-primary {
  background-color: #bea202;
  border-color: #bea202;
}
a.button.button-primary:hover, a.button.button-primary:active, a.button.button-primary:visited {
  background-color: #dbbb01;
  border-color: #bea202;
}
.signup {
  margin-right: 2rem;
}
</style>
