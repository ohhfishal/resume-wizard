<script lang="ts">
	import Header from '$lib/components/Header.svelte';
	import type { Resume } from '$lib/types';

	let resumes: Resume[] = [
		{ id: 1, title: 'Software Engineer', created: new Date('September 17, 2025') }
	];

	function addNewRow(): void {
		resumes = [
			...resumes,
			{ id: resumes.length + 1, title: 'TODO: Implement', created: new Date() }
		];
		console.log('CALLED');
	}
</script>

<Header />
<div class="container">
	<div class="hero">
		<h1>Home</h1>
		<p>Put a search bar or something here.</p>
	</div>
	<div class="card">
		<div class="card-header">
			<h3>Base Resumes</h3>
			<div class="card-meta">
				<span class="count-badge">{resumes.length} Resumes</span>
				<a class="button btn-add" href="/resume">+ Add New </a>
			</div>
		</div>
		<div class="resume-list">
			{#each resumes as resume (resume.id)}
				<div class="resume-item">
					<div class="resume-info">
						<h4>{resume.title}</h4>
						<p>Applied {resume.created.toDateString()}</p>
					</div>
					<div class="resume-actions">
						<a href="/tailor?id={resume.id}" class="button btn-small">Tailor</a>
						<a href="/view/base?id=1" class="button btn-small">View</a>
						<a href="/edit" class="button btn-small">Edit</a>
					</div>
				</div>
			{/each}
		</div>
	</div>
</div>
