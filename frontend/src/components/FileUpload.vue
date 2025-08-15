<template>
  <div class="w-full max-w-4xl mx-auto p-6">
    <!-- Header -->
    <div class="relative text-center mb-8">
      <!-- History Button -->
      <button
        @click="showHistory = true"
        class="absolute top-0 right-0 p-2 text-gray-600 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
        title="Upload History"
      >
        <svg
          class="w-6 h-6"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
          ></path>
        </svg>
      </button>

      <h1 class="text-3xl font-bold text-gray-900 mb-2">
        📁 File Upload to Telegram
      </h1>
      <p class="text-gray-600">Secure, Fast, Permanent Storage</p>
    </div>

    <!-- Upload Area -->
    <Card class="p-6 mb-6">
      <div class="flex flex-col lg:flex-row gap-6 items-center overflow-hidden">
        <!-- Upload Zone -->
        <div
          ref="dropZone"
          :class="[
            'border-2 border-dashed rounded-lg p-8 text-center transition-all duration-800 flex-1 lg:flex-initial',
            isDragOver
              ? 'border-green-400 bg-green-50 scale-105'
              : 'border-blue-400 bg-blue-50 hover:border-blue-500 hover:bg-blue-100',
          ]"
          style="min-width: 280px; flex: 1 1 500px; max-width: none"
          @dragenter.prevent="handleDragEnter"
          @dragover.prevent="handleDragOver"
          @dragleave.prevent="handleDragLeave"
          @drop.prevent="handleDrop"
        >
          <div class="text-4xl mb-4">📤</div>
          <p class="text-gray-700 mb-4">Drag files here or click to select</p>

          <input
            ref="fileInput"
            type="file"
            multiple
            class="hidden"
            @change="handleFileSelect"
          />

          <Button
            @click="() => fileInput?.click()"
            :variant="selectedFiles.length > 0 ? 'secondary' : 'default'"
            class="mb-4"
          >
            {{
              selectedFiles.length > 0
                ? `${selectedFiles.length} file(s) selected`
                : "Choose Files"
            }}
          </Button>

          <div class="flex justify-center gap-4 text-sm text-gray-500">
            <span>📊 Multiple Files</span>
            <span>🔒 Secure</span>
            <span>⚡ Fast</span>
          </div>
        </div>

        <div class="flex-shrink-0">
          <Cat />
        </div>
      </div>

      <!-- Share to Plaza Option -->
      <div
        v-if="selectedFiles.length > 0"
        class="flex items-center justify-center gap-2 mt-4"
      >
        <input
          id="shareToPlaza"
          type="checkbox"
          v-model="shareToPlaza"
          class="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500"
        />
        <label for="shareToPlaza" class="text-sm text-gray-700 cursor-pointer">
          🌍 Share to Plaza (Public)
        </label>
      </div>

      <Button
        @click="startUpload"
        :disabled="selectedFiles.length === 0 || isUploading"
        class="w-full mt-4"
        size="lg"
      >
        {{
          isUploading
            ? "Uploading..."
            : selectedFiles.length === 0
              ? "🈳 No file"
              : "⬆️ UPLOAD"
        }}
      </Button>
    </Card>

    <!-- Upload Progress -->
    <div v-if="uploadTasks.length > 0" class="space-y-4 mb-6">
      <div
        v-for="task in uploadTasks"
        :key="task.id"
        class="bg-white rounded-lg border p-4 shadow-sm"
      >
        <div class="flex justify-between items-center mb-2">
          <span class="font-medium text-gray-900 truncate flex-1 mr-4">
            {{ task.fileName }}
          </span>
          <span class="text-sm text-gray-500 whitespace-nowrap">
            {{ Math.round(task.progress) }}%
          </span>
        </div>

        <div class="mb-2">
          <Progress :model-value="task.progress" class="h-2" />
        </div>

        <div class="flex justify-between items-center text-xs text-gray-500">
          <span>{{ task.status }}</span>
          <span v-if="task.speed && task.eta">
            {{ formatSpeed(task.speed) }} • ETA: {{ formatTime(task.eta) }}
          </span>
        </div>
      </div>
    </div>

    <!-- Upload Results -->
    <div v-if="uploadResults.length > 0" class="space-y-4">
      <div
        v-for="result in uploadResults"
        :key="result.id"
        :class="[
          'rounded-lg border-l-4 p-4 shadow-sm',
          result.success
            ? 'border-l-green-500 bg-green-50'
            : 'border-l-red-500 bg-red-50',
        ]"
      >
        <div class="flex items-start justify-between">
          <div class="flex-1">
            <div class="font-medium mb-2">
              {{ result.success ? "✅" : "❌" }} {{ result.fileName }}
            </div>

            <div
              v-if="result.success && result.data?.shortFileUrl"
              class="mb-3"
            >
              <a
                :href="result.data.shortFileUrl"
                target="_blank"
                class="text-blue-600 hover:text-blue-800 underline break-all"
              >
                {{ result.data.shortFileUrl }}
              </a>
            </div>

            <div
              v-if="result.success && result.data?.shortFileUrl"
              class="flex flex-wrap gap-2"
            >
              <Button
                size="sm"
                variant="outline"
                @click="
                  (e: Event) =>
                    copyToClipboard(
                      result?.data?.shortFileUrl || '',
                      'Link copied!',
                      e
                    )
                "
              >
                Copy Link
              </Button>

              <Button
                v-if="isImageFile(result.fileName)"
                size="sm"
                variant="outline"
                @click="
                  (e: Event) =>
                    copyHtmlCode(
                      result?.data?.shortFileUrl || '',
                      result.fileName,
                      e
                    )
                "
              >
                HTML
              </Button>

              <Button
                v-if="isImageFile(result.fileName)"
                size="sm"
                variant="outline"
                @click="
                  (e: Event) =>
                    copyMarkdownCode(
                      result?.data?.shortFileUrl || '',
                      result.fileName,
                      e
                    )
                "
              >
                Markdown
              </Button>
            </div>

            <div v-if="!result.success" class="text-red-600">
              {{ result.error }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- History Modal -->
    <HistoryModal :is-open="showHistory" @close="showHistory = false" />

    <!-- Comments Section -->
    <div class="mt-12 pt-8 border-t border-gray-200">
      <h2 class="text-xl font-semibold text-gray-900 mb-6 text-center">
        💬 Comments & Feedback
      </h2>
      <GiscusComments />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import {
  uploadFile,
  uploadChunk,
  mergeChunks,
  type UploadResponse,
} from "@/services/api";
import {
  formatSpeed,
  formatTime,
  generateUploadId,
  copyToClipboard as copyText,
} from "@/lib/utils";
import { getUserFingerprintId } from "@/lib/fingerprint";
import Button from "@/components/ui/Button.vue";
import Card from "@/components/ui/Card.vue";
import Progress from "@/components/ui/Progress.vue";
import HistoryModal from "@/components/HistoryModal.vue";
import GiscusComments from "@/components/GiscusComments.vue";
import Cat from "@/components/Cat.vue";

const CHUNK_SIZE = 10 * 1024 * 1024; // 10MB
const MAX_CONCURRENT_UPLOADS = 3;

interface UploadTask {
  id: string;
  fileName: string;
  progress: number;
  status: string;
  speed?: number;
  eta?: number;
}

interface UploadResult {
  id: string;
  fileName: string;
  success: boolean;
  data?: UploadResponse;
  error?: string;
}

const dropZone = ref<HTMLDivElement>();
const fileInput = ref<HTMLInputElement>();
const selectedFiles = ref<File[]>([]);
const isDragOver = ref(false);
const isUploading = ref(false);
const uploadTasks = ref<UploadTask[]>([]);
const uploadResults = ref<UploadResult[]>([]);
const showHistory = ref(false);
const userFingerprint = ref<string>("");
const shareToPlaza = ref(false);

const handleDragEnter = (e: DragEvent) => {
  e.preventDefault();
  isDragOver.value = true;
};

const handleDragOver = (e: DragEvent) => {
  e.preventDefault();
  isDragOver.value = true;
};

const handleDragLeave = (e: DragEvent) => {
  e.preventDefault();
  if (!dropZone.value?.contains(e.relatedTarget as Node)) {
    isDragOver.value = false;
  }
};

const handleDrop = (e: DragEvent) => {
  e.preventDefault();
  isDragOver.value = false;

  const files = Array.from(e.dataTransfer?.files || []);
  selectedFiles.value = files;
};

const handleFileSelect = (e: Event) => {
  const target = e.target as HTMLInputElement;
  const files = Array.from(target.files || []);
  selectedFiles.value = files;
};

const isImageFile = (fileName: string): boolean => {
  const imageExts = [".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".svg"];
  const ext = fileName.toLowerCase().substring(fileName.lastIndexOf("."));
  return imageExts.includes(ext);
};

const copyToClipboard = async (
  text: string,
  successMessage: string,
  event?: Event
) => {
  const success = await copyText(text);
  if (success && event?.target) {
    // Store original button text
    const button = event.target as HTMLButtonElement;
    const originalText = button.textContent;
    const originalBgColor = button.style.backgroundColor;
    const originalColor = button.style.color;
    button.style.backgroundColor = "#10b981";
    button.style.color = "white";
    // Show success message
    button.textContent = successMessage;
    // Restore original text after 1 second
    setTimeout(() => {
      button.textContent = originalText;
      button.style.backgroundColor = originalBgColor;
      button.style.color = originalColor;
    }, 1000);
  } else {
    console.error("Copy faild");
  }
};

const copyHtmlCode = async (url: string, fileName: string, event?: Event) => {
  const htmlCode = `<img src="${url}" alt="${fileName}">`;
  await copyToClipboard(htmlCode, "HTML copied!", event);
};

const copyMarkdownCode = async (
  url: string,
  fileName: string,
  event?: Event
) => {
  const markdownCode = `![${fileName}](${url})`;
  await copyToClipboard(markdownCode, "Markdown copied!", event);
};

const createUploadTask = (fileName: string): UploadTask => {
  return {
    id: generateUploadId(),
    fileName,
    progress: 0,
    status: "Starting upload...",
    speed: 0,
    eta: 0,
  };
};

const updateTaskProgress = (
  taskId: string,
  progress: number,
  status?: string,
  speed?: number,
  eta?: number
) => {
  const task = uploadTasks.value.find((t) => t.id === taskId);
  if (task) {
    task.progress = progress;
    if (status) task.status = status;
    if (speed !== undefined) task.speed = speed;
    if (eta !== undefined) task.eta = eta;
  }
};

const removeTask = (taskId: string) => {
  const index = uploadTasks.value.findIndex((t) => t.id === taskId);
  if (index > -1) {
    uploadTasks.value.splice(index, 1);
  }
};

const addResult = (result: UploadResult) => {
  uploadResults.value.unshift(result);
};

const uploadSmallFile = async (file: File): Promise<void> => {
  const task = createUploadTask(file.name);
  uploadTasks.value.push(task);

  const startTime = Date.now();
  let lastProgress = 0;
  let lastTime = startTime;

  try {
    const result = await uploadFile(
      file,
      userFingerprint.value,
      shareToPlaza.value,
      (progress) => {
        const now = Date.now();
        const timeDiff = (now - lastTime) / 1000;
        const progressDiff = progress - lastProgress;

        if (timeDiff > 0.5 && progressDiff > 0) {
          const speed = ((progressDiff / 100) * file.size) / timeDiff;
          const remainingProgress = 100 - progress;
          const eta =
            remainingProgress > 0
              ? ((remainingProgress / 100) * file.size) / speed
              : 0;

          updateTaskProgress(task.id, progress, "Uploading...", speed, eta);
          lastProgress = progress;
          lastTime = now;
        } else {
          updateTaskProgress(task.id, progress, "Uploading...");
        }
      }
    );

    if (result.code === 0) {
      updateTaskProgress(task.id, 100, "Upload complete!");
      setTimeout(() => {
        removeTask(task.id);
        addResult({
          id: generateUploadId(),
          fileName: file.name,
          success: true,
          data: result,
        });
      }, 500);
    } else {
      throw new Error(result.message || "Upload failed");
    }
  } catch (error) {
    removeTask(task.id);
    addResult({
      id: generateUploadId(),
      fileName: file.name,
      success: false,
      error: error instanceof Error ? error.message : "Upload failed",
    });
  }
};

const uploadLargeFile = async (file: File): Promise<void> => {
  const totalChunks = Math.ceil(file.size / CHUNK_SIZE);
  const uploadId = generateUploadId();
  const chunkIds: string[] = new Array(totalChunks);

  const task = createUploadTask(file.name);
  uploadTasks.value.push(task);

  const chunkProgress = new Array(totalChunks).fill(0);
  const startTime = Date.now();

  const updateOverallProgress = () => {
    const totalProgress = chunkProgress.reduce(
      (sum, progress) => sum + progress,
      0
    );
    const overallProgress = totalProgress / totalChunks;

    const now = Date.now();
    const elapsed = (now - startTime) / 1000;
    const speed = ((overallProgress / 100) * file.size) / elapsed;
    const remainingProgress = 100 - overallProgress;
    const eta =
      remainingProgress > 0
        ? ((remainingProgress / 100) * file.size) / speed
        : 0;

    updateTaskProgress(task.id, overallProgress, "Uploading...", speed, eta);
  };

  try {
    // Upload chunks concurrently with semaphore
    const semaphore = new Semaphore(MAX_CONCURRENT_UPLOADS);
    const uploadPromises = [];

    for (let chunkIndex = 0; chunkIndex < totalChunks; chunkIndex++) {
      const promise = semaphore.acquire().then(async (release) => {
        try {
          const start = chunkIndex * CHUNK_SIZE;
          const end = Math.min(start + CHUNK_SIZE, file.size);
          const chunk = file.slice(start, end);

          const result = await uploadChunk(
            chunk,
            chunkIndex,
            uploadId,
            file.name,
            userFingerprint.value,
            (progress) => {
              chunkProgress[chunkIndex] = progress;
              updateOverallProgress();
            }
          );

          if (result.code === 0 && result.chunkId) {
            chunkIds[chunkIndex] = result.chunkId;
            chunkProgress[chunkIndex] = 100;
            updateOverallProgress();
          } else {
            throw new Error(result.message || "Chunk upload failed");
          }
        } finally {
          release();
        }
      });

      uploadPromises.push(promise);
    }

    await Promise.all(uploadPromises);

    // Merge chunks
    updateTaskProgress(task.id, 100, "Merging chunks...");

    const mergeResult = await mergeChunks({
      uploadId,
      fileName: file.name,
      chunkIds,
      fileSize: file.size,
      userFingerprint: userFingerprint.value,
      shared: shareToPlaza.value,
    });

    if (mergeResult.code === 0) {
      updateTaskProgress(task.id, 100, "Upload complete!");
      setTimeout(() => {
        removeTask(task.id);
        addResult({
          id: generateUploadId(),
          fileName: file.name,
          success: true,
          data: mergeResult,
        });
      }, 500);
    } else {
      throw new Error(mergeResult.message || "Merge failed");
    }
  } catch (error) {
    removeTask(task.id);
    addResult({
      id: generateUploadId(),
      fileName: file.name,
      success: false,
      error: error instanceof Error ? error.message : "Upload failed",
    });
  }
};

const startUpload = async () => {
  if (selectedFiles.value.length === 0 || isUploading.value) return;

  isUploading.value = true;

  try {
    for (const file of selectedFiles.value) {
      if (file.size > CHUNK_SIZE) {
        await uploadLargeFile(file);
      } else {
        await uploadSmallFile(file);
      }
    }
  } finally {
    isUploading.value = false;
    selectedFiles.value = [];
    shareToPlaza.value = false;
    if (fileInput.value) {
      fileInput.value.value = "";
    }
  }
};

// Semaphore class for controlling concurrent uploads
class Semaphore {
  private permits: number;
  private waiting: Array<() => void> = [];

  constructor(permits: number) {
    this.permits = permits;
  }

  async acquire(): Promise<() => void> {
    return new Promise((resolve) => {
      if (this.permits > 0) {
        this.permits--;
        resolve(() => this.release());
      } else {
        this.waiting.push(() => {
          this.permits--;
          resolve(() => this.release());
        });
      }
    });
  }

  private release(): void {
    this.permits++;
    if (this.waiting.length > 0) {
      const next = this.waiting.shift();
      if (next) next();
    }
  }
}

// 初始化用户指纹
onMounted(async () => {
  try {
    userFingerprint.value = await getUserFingerprintId();
    console.log("User fingerprint initialized:", userFingerprint.value);
  } catch (error) {
    console.error("Failed to initialize user fingerprint:", error);
  }
});
</script>
