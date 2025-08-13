<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    @click="closeModal"
  >
    <div
      class="bg-white rounded-lg p-6 max-w-4xl w-full mx-4 max-h-[80vh] overflow-hidden"
      @click.stop
    >
      <!-- Header -->
      <div class="flex justify-between items-center mb-4">
        <h2 class="text-xl font-bold text-gray-900">📁 File Explorer</h2>
        <button
          @click="closeModal"
          class="text-gray-500 hover:text-gray-700 text-2xl"
        >
          ×
        </button>
      </div>

      <!-- Tabs -->
      <div class="flex border-b border-gray-200 mb-4">
        <button
          @click="activeTab = 'history'"
          :class="[
            'px-4 py-2 font-medium text-sm border-b-2 transition-colors',
            activeTab === 'history'
              ? 'border-blue-500 text-blue-600'
              : 'border-transparent text-gray-500 hover:text-gray-700',
          ]"
        >
          📁 My Files
        </button>
        <button
          @click="activeTab = 'plaza'"
          :class="[
            'px-4 py-2 font-medium text-sm border-b-2 transition-colors',
            activeTab === 'plaza'
              ? 'border-blue-500 text-blue-600'
              : 'border-transparent text-gray-500 hover:text-gray-700',
          ]"
        >
          🌍 Plaza
        </button>
      </div>

      <!-- Loading State -->
      <div v-if="loading && files.length === 0" class="text-center py-8">
        <div
          class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto"
        ></div>
        <p class="mt-2 text-gray-600">
          Loading {{ activeTab === "history" ? "history" : "plaza files" }}...
        </p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="text-center py-8">
        <p class="text-red-600">{{ error }}</p>
        <Button @click="() => loadData(1)" class="mt-4">Retry</Button>
      </div>

      <!-- Empty State -->
      <div v-else-if="files.length === 0" class="text-center py-8">
        <div class="text-4xl mb-4">📂</div>
        <p class="text-gray-600">
          {{
            activeTab === "history"
              ? "No upload history found"
              : "No shared files in plaza"
          }}
        </p>
      </div>

      <!-- Files List -->
      <div v-else class="overflow-y-auto max-h-96">
        <div class="space-y-3">
          <div
            v-for="file in files"
            :key="file.fileId"
            class="border rounded-lg p-4 hover:bg-gray-50 transition-colors"
          >
            <div class="flex justify-between items-start">
              <div class="flex-1">
                <div class="flex items-center gap-2">
                  <h3 class="font-medium text-gray-900 truncate">
                    {{ file.filename }}
                  </h3>
                  <span
                    v-if="file.shared"
                    class="text-xs bg-green-100 text-green-800 px-2 py-1 rounded-full"
                  >
                    🌍 Shared
                  </span>
                </div>
                <p class="text-sm text-gray-500 mt-1">
                  {{ formatDate(file.time) }}
                </p>
              </div>
              <div class="flex gap-2 ml-4">
                <Button
                  size="sm"
                  variant="outline"
                  @click="
                    (e: Event) => copyFileLink(file.fileId || '', 'Copied!', e)
                  "
                >
                  Copy Link
                </Button>
                <Button
                  size="sm"
                  variant="outline"
                  @click="openFile(file.fileId)"
                >
                  Open
                </Button>
              </div>
            </div>
          </div>
        </div>

        <!-- Load More Button -->
        <div v-if="pagination.hasMore" class="text-center mt-4">
          <Button
            @click="loadMore"
            :disabled="loading"
            variant="outline"
            class="w-full"
          >
            {{ loading ? "Loading..." : "Load More" }}
          </Button>
        </div>

        <!-- Pagination Info -->
        <div
          v-if="files.length > 0"
          class="text-center mt-4 text-sm text-gray-500"
        >
          Showing {{ files.length }} of {{ pagination.total }} files
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from "vue";
import {
  getUserHistory,
  getPlazaFiles,
  type FileRecord,
  type PaginationInfo,
} from "@/services/api";
import { getUserFingerprintId } from "@/lib/fingerprint";
import { copyToClipboard } from "@/lib/utils";
import Button from "@/components/ui/Button.vue";

interface Props {
  isOpen: boolean;
}

interface Emits {
  (e: "close"): void;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

const loading = ref(false);
const error = ref("");
const files = ref<FileRecord[]>([]);
const activeTab = ref<"history" | "plaza">("history");
const pagination = ref<PaginationInfo>({
  page: 1,
  pageSize: 20,
  total: 0,
  hasMore: false,
});

const closeModal = () => {
  emit("close");
};

const loadData = async (page = 1) => {
  loading.value = true;
  error.value = "";

  try {
    let response;

    if (activeTab.value === "history") {
      const fingerprint = await getUserFingerprintId();
      response = await getUserHistory(
        fingerprint,
        page,
        pagination.value.pageSize
      );
    } else {
      response = await getPlazaFiles(page, pagination.value.pageSize);
    }

    if (response.code === 0 && response.data) {
      const plazzFiles = response.data.files || [];
      if (page === 1) {
        files.value = plazzFiles;
      } else {
        files.value.push(...plazzFiles);
      }
      pagination.value = response.data.pagination;
    } else {
      error.value = response.message || "Failed to load data";
    }
  } catch (err) {
    error.value = "Network error occurred";
    console.error("Failed to load data:", err);
  } finally {
    loading.value = false;
  }
};

const loadMore = () => {
  if (pagination.value.hasMore && !loading.value) {
    loadData(pagination.value.page + 1);
  }
};

const formatDate = (dateString: string) => {
  const date = new Date(dateString);
  return date.toLocaleString();
};

const copyFileLink = async (
  fileId: string,
  successMessage: string,
  event: Event
) => {
  const link = `${window.location.origin}/d/${fileId}`;
  const success = await copyToClipboard(link);
  if (success) {
    const button = event.target as HTMLButtonElement;
    const originalText = button.textContent;
    const originalBgColor = button.style.backgroundColor;
    const originalColor = button.style.color;
    button.style.backgroundColor = "#10b981";
    button.style.color = "white";
    button.textContent = successMessage;
    setTimeout(() => {
      button.textContent = originalText;
      button.style.backgroundColor = originalBgColor;
      button.style.color = originalColor;
    }, 1000);
  }
};

const openFile = (fileId: string) => {
  const link = `${window.location.origin}/d/${fileId}`;
  window.open(link, "_blank");
};

// Watch for tab changes
watch(activeTab, () => {
  files.value = [];
  pagination.value.page = 1;
  loadData(1);
});

// Load data when modal opens
watch(
  () => props.isOpen,
  (isOpen) => {
    if (isOpen) {
      loadData(1);
    }
  }
);
</script>
