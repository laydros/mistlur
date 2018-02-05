package main

import (
    "github.com/murlokswarm/app"
    "math"
    "math/cmplx"
    "time"
)

// Player is the component displaying Player.
type Player struct {
    Time int
    Bar  [10]float64
}

type Tag struct {
    Artist string
    Title  string
}

var (
    tag     Tag
    guidone chan struct{}
)

func (p *Player) OnMount() {
    go func() {
        c := time.Tick(60 * time.Millisecond)
        for _ = range c {

            select {
            default:
                app.Render(p)

                for i := 0; i < FFTSamples; i++ {
                    csamples[i] = complex((samples[i][0] + samples[i][1]), 0)
                }

                fftc.Transform(csamples)

                for j := 0; j < len(p.Bar); j++ {
                    for i := 0; i < len(csamples)/2/len(p.Bar); i++ {
                        p.Bar[j] = 20 * (math.Log(1 + cmplx.Abs(csamples[i+j])))
                    }
                }

            case <-guidone:
                return
            }
        }
    }()
}

func (p *Player) Next() {
    done <- struct{}{}
}

func (p *Player) OnDismount() {
    guidone <- struct{}{}
    done <- struct{}{}
}

func (p *Player) Render() string {
    return `
<div class="center">
    <div class="graph">
            <div style="height: 100px; background-color: rgba(0,0,0,0)" class="bar"/>
                {{ range $key, $data := .Bar }}
                    <div style="height: {{$data}}px;" class="bar"/>
                {{ end }}
            <div style="height: 100px; background-color: rgba(0,0,0,0)" class="bar"/>
    </div>
</div>
<h1>` + tag.Artist + `</h1>
<h2>` + tag.Title + `</h2>
<div>
    <button class="button back" onclick="OK"/>
    <button class="button play" onclick="TogglePlay"/>
    <button class="button next" onclick="Next"/>                
</div>
`
}

// /!\ Register the component. Required to use the component into a context.
func init() {
    app.RegisterComponent(&Player{})
}
