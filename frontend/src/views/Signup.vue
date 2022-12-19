<script setup></script>

<template>
  <div class="container">
    <div class="row">
      <form onsubmit="return false;">
      <div class="row">
        <div class="six columns">
          <label for="exampleEmailInput">Your email</label>
          <input class="u-full-width" type="email" placeholder="test@mailbox.com" id="email">
        </div>
      </div>
      <div class="row">
        <div class="six columns">
          <label for="exampleEmailInput">Your password</label>
          <input class="u-full-width" type="password" id="password">
        </div>
      </div>
      <input class="button-primary" type="submit" value="Submit">
      </form>
    </div>
  </div>

</template>

<style scoped>
.row {
  color: #e0bf0e;
}

input[type="submit"].button-primary, input[type="submit"].button-primary:hover {
  background-color: #bea202;
  border-color: #bea202;
}
</style>