<template>
  <div>
    <h1>Study Activity Popup</h1>
    <slot></slot>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';

export default defineComponent({
  name: 'StudyActivityPopup',
  props: {
    activityId: {
      type: Number,
      required: true
    },
    initiallyOpen: {
      type: Boolean,
      required: true
    },
    onClose: {
      type: Function,
      default: () => {}
    }
  }
});
</script>

<style scoped>
/* Add your styles here */
</style>
