#!/usr/bin/python

import sys
import pygame

pygame.init()


class GameLauncher():
    def __init__(self, screen):
        self.screen = screen
        self.scr_width = self.screen.get_rect().width
        self.scr_height = self.screen.get_rect().height

        self.clock = pygame.time.Clock()

        self.offset = 0
        self.speed = 0
        self.next_stop = None
        self.cur_pos = 0

        self._load_games_data()

        for index, game in enumerate(self.games):
            game['screenImg'] = pygame.image.load('assets/' + game['screenFile'])

        self.font = pygame.font.Font('assets/SourceSansPro-Regular.ttf', 30)

    def _load_games_data(self):
        self.games = []

        with open('assets/games.csv') as f:
            for line in f.readlines():
                data = line.split(',')
                game = {'pkgName': data[0],
                        'name': data[1],
                        'screenFile': data[2],
                        'description': data[3]}
                self.games.append(game)

    def _is_up(self, keys):
        return keys[pygame.K_UP] or keys[pygame.K_w]

    def _is_down(self, keys):
        return keys[pygame.K_DOWN] or keys[pygame.K_s]

    def _compute_offset(self, dt):
        pygame.event.pump()
        keys = pygame.key.get_pressed()

        if self._is_up(keys) and self.cur_pos > 0:
            self.speed = 200
            if self.next_stop == None:
                self.next_stop = self.offset + 20
                self.cur_pos = self.cur_pos - 1

        elif self._is_down(keys) and self.cur_pos < len(self.games) - 1:
            self.speed = -200
            if self.next_stop == None:
                self.next_stop = self.offset - 20
                self.cur_pos = self.cur_pos + 1

        if self.speed != 0:
            self.offset = self.offset + (self.speed * dt)
            if self.speed > 0 and self.offset >= self.next_stop:
                self.offset = self.next_stop
                self.speed = 0
                self.next_stop = None
            elif self.speed < 0 and self.offset <= self.next_stop:
                self.offset = self.next_stop
                self.speed = 0
                self.next_stop = None

    def run(self):
        mainloop = True
        while mainloop:
            # Limit frame speed to 50 FPS
            self.clock.tick(50)

            for event in pygame.event.get():
                if event.type == pygame.QUIT:
                    mainloop = False
                if event.type == pygame.KEYDOWN:
                    if event.key == pygame.K_ESCAPE:
                        mainloop = False
            self._compute_offset(1.0/50.0)


            # Redraw the background
            self.screen.fill((0, 0, 0))  # black

            for index, game in enumerate(self.games):
                if index == self.cur_pos:
                    label = self.font.render(game['name'], True, (255, 50, 50))
                else:
                    label = self.font.render(game['name'], True, (100, 255, 100))

                pos_x = self.scr_width / 20
                pos_y = self.scr_height / 20 + index * 30 + self.offset
                self.screen.blit(label, (pos_x, pos_y))

            img = self.games[self.cur_pos]['screenImg']
            self.screen.blit(img, (5 * self.scr_width / 20, self.scr_height / 20))
            pygame.display.flip()

if __name__ == "__main__":
    # Creating the screen
    screen = pygame.display.set_mode((1400, 800), 0, 32)

    pygame.display.set_caption('Game Launcher')
    gl = GameLauncher(screen)
    gl.run()
