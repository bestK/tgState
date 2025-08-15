<template>
  <div class="cat-video">
    <video
      ref="videoRef"
      autoplay
      muted
      loop
      playsinline
      :class="videoClass"
      @timeupdate="handleTimeUpdate"
    >
      <source
        src="http://tgstate.linkof.link/d/BQACAgUAAyEGAASBAAHQNwACATtonvW_YOvmnmAoCXwAAfvOY9ZY1w8AAocXAAJpcPlUfl4m7iXOrpE2BA"
        type="video/mp4"
      />
      Your browser does not support the video tag.
    </video>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";

interface Props {
  width?: string;
  maxWidth?: string;
  borderRadius?: string;
}

const props = withDefaults(defineProps<Props>(), {
  width: "100%",
  maxWidth: "280px",
  borderRadius: "8px",
});

const videoRef = ref<HTMLVideoElement>();
const isExpanded = ref(false);
const isHidden = ref(false);

const videoClass = computed(() => [
  "cat-video-element",
  {
    expanded: isExpanded.value,
    hidden: isHidden.value,
  },
]);

const handleTimeUpdate = () => {
  if (!videoRef.value) return;

  const currentTime = videoRef.value.currentTime;
  const duration = videoRef.value.duration;

  // 在5-10秒之间扩大视频
  if (currentTime >= 5 && currentTime <= 10) {
    isExpanded.value = true;
    isHidden.value = false;
  } else if (duration && currentTime >= duration - 2) {
    // 视频结束前2秒开始隐藏
    isExpanded.value = false;
    isHidden.value = true;
  } else {
    isExpanded.value = false;
    isHidden.value = false;
  }
};

onMounted(() => {
  // 确保视频能够正常播放
  if (videoRef.value) {
    videoRef.value.play().catch(console.error);
  }
});
</script>

<style scoped>
.cat-video {
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
}

.cat-video-element {
  width: v-bind("props.width");
  max-width: v-bind("props.maxWidth");
  height: auto;
  border: none;
  outline: none;
  display: block;
  border-radius: v-bind("props.borderRadius");
  transition: all 0.8s cubic-bezier(0.4, 0, 0.2, 1);
  transform-origin: center;
}

.cat-video-element.expanded {
  max-width: 450px;
  transform: scale(1.1);
  z-index: 10;
  position: relative;
}

.cat-video-element.hidden {
  max-width: 0;
  width: 0;
  opacity: 0;
  transform: scale(0);
}
</style>
