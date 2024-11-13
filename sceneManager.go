package main

type SceneName string

type Scene interface {
	UpdateScene(sceneManager *SceneManager)
}

type SceneManager struct {
	CurrentScene Scene
}

func NewSceneManager() *SceneManager {
	return &SceneManager{}
}

func (sm *SceneManager) SetScene(scene Scene) {
	sm.CurrentScene = scene
}

func (sm *SceneManager) Update() {
	if sm.CurrentScene != nil {
		sm.CurrentScene.UpdateScene(sm)
	}
}
