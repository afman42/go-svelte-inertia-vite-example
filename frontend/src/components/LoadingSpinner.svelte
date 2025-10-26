<script lang="ts">
  let { size = 'md', color = 'purple', centered = false }: { size?: 'sm' | 'md' | 'lg', color?: string, centered?: boolean } = $props();
  
  let sizeClass = $derived(size === 'sm' ? '4' : size === 'lg' ? '8' : '6');
  let colorClasses = $derived(`border-t-2 border-${color}-500 border-r-2 border-${color}-500`);
</script>

<div class:centered={centered ? 'flex items-center justify-center' : ''}>
  <div class={`animate-spin rounded-full h-${sizeClass} w-${sizeClass} ${colorClasses}`}></div>
</div>