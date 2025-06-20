<script lang="ts">
  import { Dialog } from "bits-ui";
  import type { Snippet } from "svelte";
  import type { HTMLAttributes } from "svelte/elements";
  
  type ModalSize = "sm" | "default" | "lg" | "xl" | "full";
  
  interface ModalProps extends HTMLAttributes<HTMLDivElement> {
    open?: boolean;
    title?: string;
    description?: string;
    closeOnOutsideClick?: boolean;
    closeOnEscape?: boolean;
    size?: ModalSize;
    children?: Snippet;
    footer?: Snippet;
  }
  
  let {
    open = $bindable(false),
    title = "",
    description = "",
    size = "default",
    children,
    footer,
    ...restProps
  }: ModalProps = $props();
  
  // Size variants
  const sizeClasses: Record<ModalSize, string> = {
    sm: "max-w-md",
    default: "max-w-lg", 
    lg: "max-w-2xl",
    xl: "max-w-4xl",
    full: "max-w-[95vw] max-h-[95vh]"
  };
</script>

<Dialog.Root bind:open>  
  <Dialog.Portal>
    <Dialog.Overlay 
      class="fixed inset-0 z-50 bg-black/50 backdrop-blur-sm"
    />
    
    <Dialog.Content
      class="fixed left-1/2 top-1/2 z-50 w-full {sizeClasses[size]} -translate-x-1/2 -translate-y-1/2 transform rounded-lg border bg-white p-6 shadow-lg dark:bg-gray-800 dark:border-gray-700"
    >
      {#if title || description}
        <Dialog.Title class="mb-4">
          {#if title}
            <Dialog.Title class="text-lg font-semibold text-gray-900 dark:text-gray-100">
              {title}
            </Dialog.Title>
          {/if}
          {#if description}
            <Dialog.Description class="mt-2 text-sm text-gray-600 dark:text-gray-400">
              {description}
            </Dialog.Description>
          {/if}
        </Dialog.Title>
      {/if}
      
      <Dialog.Close
        class="absolute right-4 top-4 rounded-sm opacity-70 ring-offset-white transition-opacity hover:opacity-100 focus:outline-none focus:ring-2 focus:ring-gray-400 focus:ring-offset-2 dark:ring-offset-gray-800 dark:focus:ring-gray-600"
      >
        <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
        <span class="sr-only">Close</span>
      </Dialog.Close>
      
      <div class="modal-content">
        {#if children}
          {@render children()}
        {/if}
      </div>
    </Dialog.Content>
  </Dialog.Portal>
</Dialog.Root>