<template>
  <div class="preview-container">
    <div class="pane-header">プレビュー</div>
    <div class="preview-content" v-html="renderedHtml"></div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  renderedHtml?: string
}

const props = defineProps<Props>()

const renderedHtml = computed(() => {
  if (props.renderedHtml && props.renderedHtml.trim()) {
    return props.renderedHtml
  }
  return '<p style="color: #6a737d; font-style: italic;">プレビューがここに表示されます...</p>'
})
</script>

<style scoped>
.preview-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  background-color: #fafafa;
}

.pane-header {
  padding: 8px 12px;
  background-color: #e9ecef;
  border-bottom: 1px solid #ddd;
  font-weight: bold;
  font-size: 14px;
}

.preview-content {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  background-color: #fff;
  margin: 10px;
  border-radius: 4px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

/* Markdownスタイル */
.preview-content :deep(h1),
.preview-content :deep(h2),
.preview-content :deep(h3),
.preview-content :deep(h4),
.preview-content :deep(h5),
.preview-content :deep(h6) {
  margin-top: 24px;
  margin-bottom: 12px;
  line-height: 1.25;
}

.preview-content :deep(h1) {
  font-size: 2em;
}

.preview-content :deep(h2) {
  font-size: 1.5em;
}

.preview-content :deep(h3) {
  font-size: 1.25em;
}

.preview-content :deep(p) {
  margin-bottom: 16px;
  line-height: 1.6;
}

.preview-content :deep(code) {
  background-color: #f6f8fa;
  padding: 2px 4px;
  border-radius: 3px;
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 0.9em;
}

.preview-content :deep(pre) {
  background-color: #f6f8fa;
  padding: 16px;
  border-radius: 6px;
  overflow-x: auto;
  line-height: 1.4;
}

.preview-content :deep(pre code) {
  background-color: transparent;
  padding: 0;
}

.preview-content :deep(blockquote) {
  border-left: 4px solid #dfe2e5;
  padding-left: 16px;
  margin: 16px 0;
  color: #6a737d;
}

.preview-content :deep(table) {
  border-collapse: collapse;
  width: 100%;
  margin: 16px 0;
}

.preview-content :deep(th),
.preview-content :deep(td) {
  border: 1px solid #dfe2e5;
  padding: 8px 12px;
  text-align: left;
}

.preview-content :deep(th) {
  background-color: #f6f8fa;
  font-weight: 600;
}

.preview-content :deep(ul),
.preview-content :deep(ol) {
  padding-left: 20px;
  margin-bottom: 16px;
}

.preview-content :deep(li) {
  margin-bottom: 8px;
}

/* レスポンシブデザイン */
@media (max-width: 768px) {
  .preview-container {
    max-height: 50vh;
  }
}
</style>
