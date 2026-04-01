<script setup lang="ts">
import type { PropType } from 'vue'

const props = defineProps({
  icon: {
    type: [Object, Function] as PropType<any>,
    required: true,
  },
  title: {
    type: String,
    required: true,
  },
  description: {
    type: String,
    required: true,
  },
  prefix: {
    type: String,
    default: '',
  },
  cardClass: {
    type: String,
    default: '',
  },
  iconWrapClass: {
    type: String,
    default: '',
  },
  iconClass: {
    type: String,
    default: 'h-4 w-4',
  },
  titleClass: {
    type: String,
    default: 'mt-3 text-sm font-bold tracking-wider text-pt-orange',
  },
  descriptionClass: {
    type: String,
    default: 'mt-1 text-sm leading-relaxed tracking-wider text-black',
  },
})
</script>

<template>
  <article
    class="rounded-2xl border border-[#C5CBE3] bg-white p-5 shadow-sm transition"
    :class="props.cardClass"
  >
    <div
      class="flex h-9 w-full items-center justify-center !bg-pt-orange rounded-lg text-white"
      :class="props.iconWrapClass"
    >
      <component :is="props.icon" :class="props.iconClass" />
    </div>
    <h3 :class="props.titleClass">
      <template v-if="props.prefix">{{ props.prefix }}. </template>
      {{ props.title }}
    </h3>
    <p :class="props.descriptionClass">{{ props.description }}</p>
  </article>
</template>
