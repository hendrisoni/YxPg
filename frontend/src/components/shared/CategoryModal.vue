<template>
  <Modal :show="true" title="Add Category" @close="$emit('close')" size="sm">
    <div class="space-y-4 py-2 select-none">
      <div class="flex flex-col gap-2">
        <label class="block text-xs font-medium text-text-secondary">Category Name <span class="text-red-500">*</span></label>
        <input 
          ref="categoryInput" 
          v-model="categoryName" 
          type="text" 
          class="w-full bg-navy-tertiary border border-navy-border rounded px-3 py-1.5 text-sm text-text-primary focus:border-teal-accent focus:outline-none" 
          placeholder="e.g. My Project" 
          @keyup.enter="saveCategory" 
        />
      </div>
      
      <div class="pt-4 flex justify-end gap-2">
        <button 
          @click="$emit('close')" 
          class="px-3 py-1.5 text-xs font-medium text-text-secondary hover:text-text-primary hover:bg-navy-hover rounded transition-colors"
        >
          Cancel
        </button>
        <button 
          @click="saveCategory" 
          class="px-3 py-1.5 text-xs font-medium bg-teal-accent text-navy-primary hover:bg-teal-hover rounded transition-colors disabled:opacity-50 disabled:cursor-not-allowed" 
          :disabled="!categoryName.trim()"
        >
          Save
        </button>
      </div>
    </div>
  </Modal>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useWorkspaceStore } from '../../stores/workspace'
import Modal from './Modal.vue'

const emit = defineEmits(['close'])
const workspaceStore = useWorkspaceStore()
const categoryName = ref('')
const categoryInput = ref<HTMLInputElement | null>(null)

onMounted(() => {
  categoryInput.value?.focus()
})

async function saveCategory() {
  if (!categoryName.value.trim()) return
  await workspaceStore.addCategory(categoryName.value.trim())
  emit('close')
}
</script>