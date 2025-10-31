<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-2xl font-bold text-gray-800">Item Categories</h1>
        <p class="text-gray-600 mt-1">Manage inventory item categories</p>
      </div>
      <Button
        v-if="canCreateCategory"
        label="Add Category"
        icon="pi pi-plus-circle"
        raised
        @click="openDialog()"
        class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 text-sm md:text-base border-none shadow-md rounded-lg"
      />
    </div>

    <!-- Categories Grid -->
    <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-2">
      <div
        v-for="category in categories"
        :key="category.id"
        class="bg-white rounded shadow-sm border border-gray-200 p-2 hover:shadow-md transition-all"
        style="min-width: 0;"
      >
        <div class="flex items-start justify-between mb-2">
          <div class="flex items-center gap-2">
            <div
              class="w-7 h-7 rounded flex items-center justify-center"
              :style="{
                backgroundColor: category.color + '20',
                color: category.color,
              }"
            >
              <i class="pi pi-tag text-base"></i>
            </div>
            <div>
              <h3 class="font-semibold text-gray-800 text-sm">{{ category.name }}</h3>
              <p class="text-xs text-gray-500">{{ category.code }}</p>
            </div>
          </div>
          <div
            v-if="canUpdateCategory || canDeleteCategory"
            class="flex gap-1"
          >
            <Button
              v-if="canUpdateCategory"
              icon="pi pi-pencil"
              class="p-button-text p-button-sm"
              style="font-size:0.85rem;"
              @click="openDialog(category)"
            />
            <Button
              v-if="canDeleteCategory"
              icon="pi pi-trash"
              class="p-button-text p-button-danger p-button-sm"
              style="font-size:0.85rem;"
              @click="confirmDelete(category)"
            />
          </div>
        </div>

        <p class="text-xs text-gray-600 mb-2 truncate">
          {{ category.description || "No description" }}
        </p>

        <div
          class="flex items-center justify-between pt-2 border-t border-gray-100"
        >
          <div class="flex items-center gap-1 text-xs text-gray-500">
            <i class="pi pi-box"></i>
            <span>{{ category.item_count || 0 }} items</span>
          </div>
          <Button
            v-if="canUpdateCategory"
            :label="category.is_active ? 'Active' : 'Inactive'"
            :icon="category.is_active ? 'pi pi-check' : 'pi pi-ban'"
            iconPos="left"
            :loading="updatingId === category.id"
            @click="toggleActive(category)"
            :class="[
              'px-3 py-2 text-sm rounded-lg transition-colors',
              category.is_active
                ? 'bg-green-100 text-green-700 hover:bg-green-200'
                : 'bg-red-100 text-red-600 hover:bg-red-200',
            ]"
          />
          <span
            v-else
            :class="[
              'px-3 py-2 text-xs font-medium rounded-lg',
              category.is_active
                ? 'bg-green-100 text-green-700'
                : 'bg-red-100 text-red-600',
            ]"
          >
            {{ category.is_active ? "Active" : "Inactive" }}
          </span>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div
      v-if="categories.length === 0"
      class="text-center py-12 bg-white rounded-lg"
    >
      <i class="pi pi-tags text-6xl text-gray-300 mb-4"></i>
      <h3 class="text-lg font-medium text-gray-800 mb-2">No Categories Yet</h3>
      <p class="text-gray-600 mb-4">Start by creating your first category</p>
      <Button
        v-if="canCreateCategory"
        label="Add Category"
        icon="pi pi-plus-circle"
        raised
        @click="openDialog()"
        class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 text-sm md:text-base border-none shadow-md rounded-lg"
      />
    </div>

    <!-- Add/Edit Dialog -->
    <Dialog
      v-model:visible="dialogVisible"
      :header="editMode ? 'Edit Category' : 'Add Category'"
      :modal="true"
      class="w-full max-w-md"
    >
      <div class="space-y-4 pt-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Category Name <span class="text-red-500">*</span>
          </label>
          <InputText
            v-model="form.name"
            placeholder="Enter category name"
            class="w-full"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Category Code <span class="text-red-500">*</span>
          </label>
          <InputText
            v-model="form.code"
            placeholder="e.g., ELEC"
            class="w-full"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Description
          </label>
          <Textarea
            v-model="form.description"
            rows="3"
            placeholder="Enter description"
            class="w-full"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Color <span class="text-red-500">*</span>
          </label>
          <div class="space-y-3">
            <div class="flex flex-wrap items-center gap-3">
              <input
                type="color"
                v-model="form.color"
                @change="form.color = normalizeHexColor(form.color) || defaultColor"
                class="w-12 h-12 border border-gray-300 rounded cursor-pointer"
              />
              <input
                v-model="form.color"
                type="text"
                maxlength="7"
                @blur="form.color = normalizeHexColor(form.color) || defaultColor"
                class="w-32 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent font-mono uppercase"
                placeholder="#FFFFFF"
              />
              <span class="text-xs text-gray-500">
                Use custom color or pick from presets below
              </span>
            </div>
            <div class="flex flex-wrap gap-2">
              <button
                v-for="color in colorPalette"
                :key="color"
                type="button"
                class="w-10 h-10 rounded-lg border-2 transition-transform"
                :class="normalizeHexColor(form.color) === normalizeHexColor(color) ? 'border-gray-900 scale-105' : 'border-transparent hover:scale-105'"
                :style="{ backgroundColor: color }"
                @click="form.color = normalizeHexColor(color) || defaultColor"
              />
            </div>
          </div>
        </div>

        <div class="flex items-center gap-2">
          <Checkbox v-model="form.is_active" inputId="is_active" binary />
          <label for="is_active" class="text-sm text-gray-700">Active</label>
        </div>
      </div>

      <template #footer>
        <Button
          label="Cancel"
          icon="pi pi-times"
          severity="danger"
          class="font-bold bg-red-500 border-none text-white hover:bg-red-600 focus:bg-red-700"
          @click="dialogVisible = false"
        />
        <Button
          :label="editMode ? 'Update' : 'Create'"
          icon="pi pi-check"
          severity="info"
          class="font-bold bg-blue-600 border-none text-white hover:bg-blue-700 focus:bg-blue-800"
          @click="saveCategory"
          :loading="loading"
        />
      </template>
    </Dialog>

    <!-- Delete Confirmation -->
    <Dialog
      v-model:visible="deleteDialogVisible"
      header="Confirm Delete"
      :modal="true"
      class="w-full max-w-md"
    >
      <div class="flex items-start gap-4">
        <i class="pi pi-exclamation-triangle text-3xl text-orange-500"></i>
        <div>
          <p class="text-gray-700">
            Are you sure you want to delete this category?
          </p>
          <p class="text-sm text-gray-500 mt-2">
            This action cannot be undone. Items in this category will need to be
            recategorized.
          </p>
        </div>
      </div>

      <template #footer>
        <Button
          label="Cancel"
          icon="pi pi-times"
          severity="danger"
          class="font-bold bg-red-500 border-none text-white hover:bg-red-600 focus:bg-red-700"
          @click="deleteDialogVisible = false"
        />
        <Button
          label="Delete"
          icon="pi pi-trash"
          severity="danger"
          class="font-bold bg-orange-500 border-none text-white hover:bg-orange-600 focus:bg-orange-700"
          @click="deleteCategory"
          :loading="loading"
        />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from "vue";
import { useToast } from "primevue/usetoast";
import axios from "@/api/axios";
import { useAuthStore } from "@/stores/auth";

const toast = useToast();
const authStore = useAuthStore();

const categories = ref([]);
const dialogVisible = ref(false);
const deleteDialogVisible = ref(false);
const loading = ref(false);
const editMode = ref(false);
const selectedCategory = ref(null);
const updatingId = ref(null);

const defaultColor = "#3B82F6";

const form = ref({
  name: "",
  code: "",
  description: "",
  color: defaultColor,
  is_active: true,
});

const colorPalette = [
  "#3B82F6", // Blue
  "#10B981", // Emerald
  "#F59E0B", // Amber
  "#EF4444", // Red
  "#8B5CF6", // Purple
  "#EC4899", // Pink
  "#06B6D4", // Cyan
  "#F97316", // Orange
  "#6366F1", // Indigo
  "#14B8A6", // Teal
  "#0EA5E9", // Sky
  "#22D3EE", // Light Cyan
  "#A855F7", // Violet
  "#F43F5E", // Rose
];

const HEX_COLOR_REGEX = /^#(?:[0-9A-Fa-f]{6}|[0-9A-Fa-f]{3})$/;

const canViewCategories = computed(() => authStore.hasPermission("category.view"));
const canCreateCategory = computed(() => authStore.hasPermission("category.create"));
const canUpdateCategory = computed(() => authStore.hasPermission("category.update"));
const canDeleteCategory = computed(() => authStore.hasPermission("category.delete"));

function normalizeHexColor(color) {
  if (typeof color !== "string") {
    return "";
  }
  let value = color.trim();
  if (!value) {
    return "";
  }
  if (!value.startsWith("#")) {
    value = `#${value}`;
  }
  if (!HEX_COLOR_REGEX.test(value)) {
    return "";
  }
  if (value.length === 4) {
    value = `#${value[1]}${value[1]}${value[2]}${value[2]}${value[3]}${value[3]}`;
  }
  return value.toUpperCase();
}

onMounted(() => {
  if (canViewCategories.value) {
    fetchCategories();
  }
});

async function fetchCategories() {
  if (!canViewCategories.value) {
    categories.value = [];
    return;
  }

  try {
    const response = await axios.get("/categories");
    categories.value = (response.data.data || []).map((category) => ({
      ...category,
      color: normalizeHexColor(category.color) || defaultColor,
    }));
  } catch (error) {
    toast.add({
      severity: "error",
      summary: "Error",
      detail: "Failed to load categories",
      life: 3000,
    });
  }
}

function openDialog(category = null) {
  if (category) {
    if (!canUpdateCategory.value) {
      toast.add({
        severity: "warn",
        summary: "Akses Ditolak",
        detail: "Anda tidak memiliki izin untuk memperbarui kategori.",
        life: 3000,
      });
      return;
    }
    editMode.value = true;
    selectedCategory.value = category;
    form.value = {
      name: category.name || "",
      code: category.code || "",
      description: category.description || "",
      color: normalizeHexColor(category.color) || defaultColor,
      is_active: Boolean(category.is_active),
    };
  } else {
    if (!canCreateCategory.value) {
      toast.add({
        severity: "warn",
        summary: "Akses Ditolak",
        detail: "Anda tidak memiliki izin untuk menambah kategori.",
        life: 3000,
      });
      return;
    }
    editMode.value = false;
    selectedCategory.value = null;
    form.value = {
      name: "",
      code: "",
      description: "",
      color: defaultColor,
      is_active: true,
    };
  }
  dialogVisible.value = true;
}

async function toggleActive(category) {
  if (!category) {
    return;
  }

  if (!canUpdateCategory.value) {
    toast.add({
      severity: "warn",
      summary: "Akses Ditolak",
      detail: "Anda tidak memiliki izin untuk memperbarui status kategori.",
      life: 3000,
    });
    return;
  }

  updatingId.value = category.id;
  const previous = category.is_active;
  const next = !previous;
  category.is_active = next;

  try {
    await axios.put(`/categories/${category.id}`, {
      is_active: next,
    });
    toast.add({
      severity: "success",
      summary: "Status Updated",
      detail: `${category.name} is now ${next ? "active" : "inactive"}`,
      life: 2500,
    });
  } catch (error) {
    category.is_active = previous;
    const message =
      error.response?.data?.message ||
      error.response?.data?.error ||
      "Failed to update status";
    toast.add({
      severity: "error",
      summary: "Update Failed",
      detail: message,
      life: 3000,
    });
  } finally {
    updatingId.value = null;
  }
}

async function saveCategory() {
  if (editMode.value && !canUpdateCategory.value) {
    toast.add({
      severity: "warn",
      summary: "Akses Ditolak",
      detail: "Anda tidak memiliki izin untuk memperbarui kategori.",
      life: 3000,
    });
    return;
  }

  if (!editMode.value && !canCreateCategory.value) {
    toast.add({
      severity: "warn",
      summary: "Akses Ditolak",
      detail: "Anda tidak memiliki izin untuk menambah kategori.",
      life: 3000,
    });
    return;
  }

  if (!form.value.name || !form.value.code) {
    toast.add({
      severity: "warn",
      summary: "Validation Error",
      detail: "Please fill in all required fields",
      life: 3000,
    });
    return;
  }

  loading.value = true;
  const normalizedColor = normalizeHexColor(form.value.color) || defaultColor;
  const payload = {
    name: form.value.name,
    code: form.value.code,
    description: form.value.description,
    color: normalizedColor,
    is_active: form.value.is_active,
  };

  try {
    if (editMode.value) {
      await axios.put(`/categories/${selectedCategory.value.id}`, payload);
      toast.add({
        severity: "success",
        summary: "Success",
        detail: "Category updated successfully",
        life: 3000,
      });
    } else {
      await axios.post("/categories", payload);
      toast.add({
        severity: "success",
        summary: "Success",
        detail: "Category created successfully",
        life: 3000,
      });
    }
    form.value.color = normalizedColor;
    dialogVisible.value = false;
    await fetchCategories();
  } catch (error) {
    toast.add({
      severity: "error",
      summary: "Error",
      detail: error.response?.data?.message || "Failed to save category",
      life: 3000,
    });
  } finally {
    loading.value = false;
  }
}

function confirmDelete(category) {
  if (!canDeleteCategory.value) {
    toast.add({
      severity: "warn",
      summary: "Akses Ditolak",
      detail: "Anda tidak memiliki izin untuk menghapus kategori.",
      life: 3000,
    });
    return;
  }

  selectedCategory.value = category;
  deleteDialogVisible.value = true;
}

async function deleteCategory() {
  if (!canDeleteCategory.value) {
    toast.add({
      severity: "warn",
      summary: "Akses Ditolak",
      detail: "Anda tidak memiliki izin untuk menghapus kategori.",
      life: 3000,
    });
    return;
  }

  loading.value = true;
  try {
    await axios.delete(`/categories/${selectedCategory.value.id}`);
    toast.add({
      severity: "success",
      summary: "Success",
      detail: "Category deleted successfully",
      life: 3000,
    });
    deleteDialogVisible.value = false;
    await fetchCategories();
  } catch (error) {
    toast.add({
      severity: "error",
      summary: "Error",
      detail: error.response?.data?.message || "Failed to delete category",
      life: 3000,
    });
  } finally {
    loading.value = false;
  }
}
</script>
