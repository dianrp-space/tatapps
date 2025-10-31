<template>
  <ul :class="listClass">
    <li
      v-for="node in nodes"
      :key="node.id"
      :class="['org-tree-item', { 'has-children': node.children && node.children.length }]"
    >
      <div class="org-node">
        <p class="org-node__title">{{ node.title }}</p>
        <p v-if="node.notes" class="org-node__notes">{{ node.notes }}</p>
      </div>
      <PositionTree
        v-if="node.children && node.children.length"
        :nodes="node.children"
      />
    </li>
  </ul>
</template>

<script setup>
import { computed } from 'vue'

defineOptions({ name: 'PositionTree' })

const props = defineProps({
  nodes: {
    type: Array,
    default: () => []
  },
  isRoot: {
    type: Boolean,
    default: false
  }
})

const listClass = computed(() => (props.isRoot ? 'org-tree-root' : 'org-tree-branch'))
</script>

<style scoped>
.org-tree-root,
.org-tree-branch {
  --node-width: 110px;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  list-style: none;
  margin: 0;
  padding: 0;
  position: relative;
  column-gap: 0;
}

.org-tree-root {
  padding-top: 0;
}

.org-tree-branch {
  padding-top: 1.1rem;
}

.org-tree-branch::before {
  content: '';
  position: absolute;
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  height: 1.1rem;
  border-left: 1px solid #d1d5db;
}

.org-tree-item {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: 1rem 0.75rem 0 0.75rem;
}

.org-tree-root > .org-tree-item {
  padding-top: 0;
}

.org-tree-branch > .org-tree-item::before,
.org-tree-branch > .org-tree-item::after {
  content: '';
  position: absolute;
  top: 0;
  width: 50%;
  height: 1.1rem;
  border-top: 1px solid #d1d5db;
}

.org-tree-branch > .org-tree-item::before {
  right: 50%;
  border-right: 1px solid #d1d5db;
}

.org-tree-branch > .org-tree-item::after {
  left: 50%;
  border-left: 1px solid #d1d5db;
}

.org-tree-branch > .org-tree-item:first-child::before {
  border-right: none;
}

.org-tree-branch > .org-tree-item:last-child::after {
  border-left: none;
}

.org-tree-branch > .org-tree-item:only-child {
  padding-top: 0;
}

.org-tree-branch > .org-tree-item:only-child::before,
.org-tree-branch > .org-tree-item:only-child::after {
  display: none;
}

.org-tree-item.has-children > .org-node::after {
  content: '';
  position: absolute;
  bottom: -1.1rem;
  left: 50%;
  transform: translateX(-50%);
  width: 0;
  height: 1.1rem;
  border-left: 1px solid #d1d5db;
}

.org-node {
  position: relative;
  width: var(--node-width);
  background: #ffffff;
  border: 1px solid #d8dee9;
  border-radius: 0.75rem;
  padding: 0.35rem 0.45rem;
  box-shadow: 0 10px 18px -28px rgba(15, 23, 42, 0.35);
  display: flex;
  flex-direction: column;
  gap: 0.15rem;
  color: #0f172a;
}

.org-node__title {
  font-weight: 600;
  font-size: 0.68rem;
  color: #0f172a;
}

.org-node__notes {
  font-size: 0.5rem;
  color: #64748b;
}

@media (max-width: 640px) {
  .org-tree-root,
  .org-tree-branch {
    flex-direction: column;
    align-items: center;
  }

  .org-tree-branch::before,
  .org-tree-branch > .org-tree-item::before,
  .org-tree-branch > .org-tree-item::after,
  .org-tree-item.has-children > .org-node::after {
    display: none;
  }

  .org-tree-root > .org-tree-item,
  .org-tree-branch > .org-tree-item {
    padding: 0.75rem 0 0 0;
  }
}
</style>
