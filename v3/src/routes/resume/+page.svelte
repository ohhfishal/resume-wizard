<script lang="ts">
	import Header from '$lib/components/Header.svelte';
  import { goto } from '$app/navigation'

  let fileInput;

  function clickFileInput(): void {
    fileInput.click();
  }

  function handleFileSelect(event): void {
    let selectedFile = event.target.files[0]
    if (selectedFile) {
      // TODO: Store the file somewhere?
      console.log("WE GOT", selectedFile.name);
      goto('/');
    }
  }

</script>

<Header />
<div class="container">
	<div class="hero">
		<h1>Create or upload a resume</h1>
		<p>You can edit it anytime.</p>
	</div>
	<div class="card">
		<div class="card-header"><h3>Pick a workflow</h3></div>
		<div class="choice-screen">
			<div class="choice-option upload" on:click={clickFileInput}>
				<div class="choice-icon">📄</div>
				<h3>Upload Resume</h3>
				<p>Have an existing resume? Upload a JSON or YAML file.</p>
				<button class="button btn-primary">Choose File</button>
        <input
          id="file-input"
          bind:this={fileInput}
          on:change={handleFileSelect}
          name="file"
          class="button btn-primary"
          type="file"
          accept=".json,.yaml"
        />
			</div>
			<a href="/base/new" class="choice-option manual"
				><div>
					<div class="choice-icon">✏️</div>
					<h3>Build from Scratch</h3>
					<p>We'll walk you through each section step by step to create your perfect resume.</p>
					<button class="button btn-primary">Get Started</button>
				</div></a
			>
		</div>
	</div>
</div>
